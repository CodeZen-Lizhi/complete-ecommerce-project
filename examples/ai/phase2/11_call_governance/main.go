package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/cloudwego/eino-ext/components/model/openai"
	einomodel "github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"golang.org/x/time/rate"
)

const (
	baseURL             = "http://localhost:8084/v1"
	apiKey              = "sk-qWpIk8nVsa8VGJyNHtcAS4VaMhCDJB1z"
	modelName           = "gpt-5.5"
	modelClientTimeout  = 30 * time.Second
	systemPrompt        = "你是一个 Go 学习助手，请用简洁、准确的中文回答。"
	maxBackoff          = 2 * time.Second
	modelCallsPerSecond = 2
	modelCallBurst      = 1
)

var (
	errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")
	modelCallLimiter      = rate.NewLimiter(rate.Limit(modelCallsPerSecond), modelCallBurst)
)

type callRequest struct {
	Prompt string
}

type callResponse struct {
	Text         string
	PromptTokens int
	OutputTokens int
}

type modelClient interface {
	Generate(ctx context.Context, request callRequest) (callResponse, error)
}

type callMetrics struct {
	Attempts int
	Latency  time.Duration
	Success  bool
	Tokens   int
}

type governanceConfig struct {
	Timeout        time.Duration
	MaxRetries     int
	InitialBackoff time.Duration
	Limiter        *rate.Limiter
}

type einoModel struct {
	chatModel einomodel.BaseChatModel
}

// Generate 通过 Eino OpenAI-compatible ChatModel 发起真实模型调用并提取 Token 用量。
func (m einoModel) Generate(ctx context.Context, request callRequest) (callResponse, error) {
	if ctx == nil {
		return callResponse{}, errors.New("Context 不能为空")
	}
	if m.chatModel == nil {
		return callResponse{}, errors.New("Eino ChatModel 未配置")
	}
	if strings.TrimSpace(request.Prompt) == "" {
		return callResponse{}, errors.New("Prompt 不能为空")
	}

	response, err := m.chatModel.Generate(ctx, []*schema.Message{
		schema.SystemMessage(systemPrompt),
		schema.UserMessage(request.Prompt),
	})
	if err != nil {
		return callResponse{}, fmt.Errorf("Eino Generate 调用失败: %w", err)
	}
	if response == nil {
		return callResponse{}, errors.New("模型响应不能为空")
	}
	if strings.TrimSpace(response.Content) == "" {
		return callResponse{}, errors.New("模型回答不能为空")
	}
	if response.ResponseMeta == nil || response.ResponseMeta.Usage == nil {
		return callResponse{}, errors.New("模型响应缺少 Token 用量")
	}
	return callResponse{
		Text:         response.Content,
		PromptTokens: response.ResponseMeta.Usage.PromptTokens,
		OutputTokens: response.ResponseMeta.Usage.CompletionTokens,
	}, nil
}

// newEinoModel 校验全局模型配置并创建可复用的真实 Eino ChatModel。
func newEinoModel(ctx context.Context) (modelClient, error) {
	if ctx == nil {
		return nil, errors.New("Context 不能为空")
	}
	if strings.TrimSpace(baseURL) == "" {
		return nil, errors.New("Base URL 未配置")
	}
	if strings.TrimSpace(apiKey) == "" || apiKey == "replace-with-your-api-key" {
		return nil, errors.New("API Key 未配置")
	}
	if strings.TrimSpace(modelName) == "" {
		return nil, errors.New("Model 未配置")
	}

	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Model:   modelName,
		Timeout: modelClientTimeout,
	})
	if err != nil {
		return nil, fmt.Errorf("创建 Eino ChatModel 失败: %w", err)
	}
	return einoModel{chatModel: chatModel}, nil
}

// main 创建真实 Model Client 并进入完整调用治理流程。
func main() {
	config := governanceConfig{
		Timeout:        modelClientTimeout,
		MaxRetries:     2,
		InitialBackoff: 200 * time.Millisecond,
		Limiter:        modelCallLimiter,
	}
	client, err := newEinoModel(context.Background())
	if err != nil {
		fmt.Printf("创建 Model Client 失败: %v\n", err)
		return
	}
	response, metrics, err := governedCall(
		context.Background(),
		client,
		config,
		callRequest{Prompt: "请解释 Go context 超时。"},
	)
	if err != nil {
		fmt.Printf("调用治理练习失败: %v\n", err)
		return
	}
	fmt.Printf("回答: %s\n尝试次数: %d\n总耗时: %s\n", response.Text, metrics.Attempts, metrics.Latency)
}

// validateGovernanceConfig 校验超时、重试次数和退避间隔。
func validateGovernanceConfig(config governanceConfig) error {
	if config.Timeout <= 0 {
		return errors.New("调用超时必须大于 0")
	}
	if config.MaxRetries < 0 {
		return errors.New("最大重试次数不能小于 0")
	}
	if config.InitialBackoff <= 0 {
		return errors.New("初始退避时间必须大于 0")
	}
	if config.Limiter == nil {
		return errors.New("调用限流器不能为空")
	}
	return nil
}

// governedCall 在超时、限流和有限重试边界内执行模型调用并记录统计。
func governedCall(
	ctx context.Context,
	client modelClient,
	config governanceConfig,
	request callRequest,
) (callResponse, callMetrics, error) {
	if ctx == nil {
		return callResponse{}, callMetrics{}, errors.New("Context 不能为空")
	}
	if client == nil {
		return callResponse{}, callMetrics{}, errors.New("Model Client 不能为空")
	}
	if err := validateGovernanceConfig(config); err != nil {
		return callResponse{}, callMetrics{}, err
	}

	startedAt := time.Now()
	metrics := callMetrics{}
	if err := config.Limiter.Wait(ctx); err != nil {
		metrics.Latency = time.Since(startedAt)
		return callResponse{}, metrics, fmt.Errorf("等待模型限流许可失败: %w", err)
	}

	var lastErr error
	// TODO 1：在外层限流器获得许可；等待时必须响应 ctx 取消。
	// TODO 2：为每次尝试创建独立超时 Context，并确保 cancel 在本轮立即执行。
	// TODO 3：分类错误，只重试临时网络错误、429 和部分 5xx；参数、认证、取消和解析错误不得重试。
	// TODO 4：使用指数退避并设置最大等待；退避期间响应 ctx 取消。
	// TODO 5：记录尝试次数、总耗时、成功状态和 Token，日志不得包含 Prompt 或 API Key。
	// TODO 6：成功返回响应和指标；耗尽重试时用 %w 保留最后一个错误。
	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		metrics.Attempts++
		attemptCtx, cancel := context.WithTimeout(ctx, config.Timeout)
		response, err := client.Generate(attemptCtx, request)
		cancel()
		metrics.Tokens += response.PromptTokens + response.OutputTokens

		if err == nil {
			metrics.Latency = time.Since(startedAt)
			metrics.Success = true
			return response, metrics, nil
		}

		lastErr = err
		if !isRetryable(ctx, err) || attempt == config.MaxRetries {
			break
		}
		if err := waitBackoff(ctx, config.InitialBackoff, attempt); err != nil {
			lastErr = fmt.Errorf("等待第 %d 次重试退避失败: %w", attempt+1, err)
			break
		}
	}

	metrics.Latency = time.Since(startedAt)
	return callResponse{}, metrics, fmt.Errorf("模型调用在 %d 次尝试后失败: %w", metrics.Attempts, lastErr)
}

// waitBackoff 执行有上限的指数退避，并在等待期间响应 Context 取消。
func waitBackoff(ctx context.Context, initial time.Duration, retryIndex int) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}
	if initial <= 0 {
		return errors.New("初始退避时间必须大于 0")
	}
	if retryIndex < 0 {
		return errors.New("重试序号不能小于 0")
	}

	delay := initial
	for i := 0; i < retryIndex && delay < maxBackoff; i++ {
		if delay > maxBackoff/2 {
			delay = maxBackoff
			break
		}
		delay *= 2
	}
	if delay > maxBackoff {
		delay = maxBackoff
	}

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// isRetryable 判断模型调用错误是否允许重试；parentCtx 必须是外层 Context。
func isRetryable(parentCtx context.Context, err error) bool {
	if err == nil {
		return false
	}

	// 整个业务请求已经取消或超时，继续重试没有意义。
	if parentCtx.Err() != nil {
		return false
	}

	// 主动取消绝对不重试。
	if errors.Is(err, context.Canceled) {
		return false
	}

	// 单次尝试超时可能已在服务端产生费用，不自动重试。
	if errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	// Eino OpenAI-compatible 服务端错误。
	var apiErr *openai.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.HTTPStatusCode {
		case http.StatusTooManyRequests, // 429
			http.StatusInternalServerError, // 500
			http.StatusBadGateway,          // 502
			http.StatusServiceUnavailable,  // 503
			http.StatusGatewayTimeout:      // 504
			return true
		default:
			// 400、401、403、404 等都不重试。
			return false
		}
	}

	// 网络超时或者明确标记为临时的网络错误。
	var netErr net.Error
	if errors.As(err, &netErr) {
		return netErr.Timeout() || netErr.Temporary()
	}

	// 参数、解析及未知错误默认不重试。
	return false
}
