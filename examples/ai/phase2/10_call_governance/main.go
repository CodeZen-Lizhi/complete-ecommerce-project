package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

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
}

type incompleteModel struct{}

// Generate 保持骨架离线，并提示学习者注入 Fake 或真实模型适配器。
func (incompleteModel) Generate(context.Context, callRequest) (callResponse, error) {
	return callResponse{}, errExerciseIncomplete
}

// main 使用离线 Model Client 进入完整调用治理流程。
func main() {
	config := governanceConfig{
		Timeout:        10 * time.Second,
		MaxRetries:     2,
		InitialBackoff: 200 * time.Millisecond,
	}
	response, metrics, err := governedCall(
		context.Background(),
		incompleteModel{},
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

	// TODO 1：在外层限流器获得许可；等待时必须响应 ctx 取消。
	// TODO 2：为每次尝试创建独立超时 Context，并确保 cancel 在本轮立即执行。
	// TODO 3：分类错误，只重试临时网络错误、429 和部分 5xx；参数、认证、取消和解析错误不得重试。
	// TODO 4：使用指数退避并设置最大等待；退避期间响应 ctx 取消。
	// TODO 5：记录尝试次数、总耗时、成功状态和 Token，日志不得包含 Prompt 或 API Key。
	// TODO 6：成功返回响应和指标；耗尽重试时用 %w 保留最后一个错误。
	return callResponse{}, callMetrics{}, errExerciseIncomplete
}
