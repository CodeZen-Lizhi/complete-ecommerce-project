package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"golang.org/x/time/rate"
)

const maximumRetryBackoff = 2 * time.Second

// governanceConfig 定义分类和回答调用共享的限流、超时与有限重试策略。
type governanceConfig struct {
	Timeout        time.Duration
	MaxRetries     int
	InitialBackoff time.Duration
	Limiter        *rate.Limiter
}

// stageMetrics 记录单个模型阶段的尝试、延迟、成功状态和可选 Token。
type stageMetrics struct {
	Attempts int
	Latency  time.Duration
	Success  bool
	Usage    tokenUsage
}

// cancelOnCloseStream 在关闭底层流后取消其专属超时 Context，避免成功流被提前取消。
type cancelOnCloseStream struct {
	stream messageStream
	cancel context.CancelFunc
	once   sync.Once
}

// Recv 将读取操作委托给底层模型流。
func (stream *cancelOnCloseStream) Recv() (modelChunk, error) {
	if stream == nil || stream.stream == nil {
		return modelChunk{}, fmt.Errorf("模型流未初始化")
	}
	return stream.stream.Recv()
}

// Close 只执行一次底层流关闭和 Context 取消，释放流式请求资源。
func (stream *cancelOnCloseStream) Close() {
	if stream == nil {
		return
	}
	stream.once.Do(func() {
		if stream.stream != nil {
			stream.stream.Close()
		}
		if stream.cancel != nil {
			stream.cancel()
		}
	})
}

// newGovernanceConfig 从练习配置创建每个进程共享的模型调用治理参数。
func newGovernanceConfig(config exerciseConfig) governanceConfig {
	return governanceConfig{
		Timeout:        config.CallTimeout,
		MaxRetries:     config.MaxRetries,
		InitialBackoff: config.InitialBackoff,
		Limiter:        rate.NewLimiter(rate.Limit(config.CallsPerSecond), 1),
	}
}

// validateGovernanceConfig 在调用 Provider 前校验治理参数。
func validateGovernanceConfig(config governanceConfig) error {
	if config.Timeout <= 0 || config.InitialBackoff <= 0 || config.MaxRetries < 0 || config.Limiter == nil {
		return fmt.Errorf("模型调用治理配置无效")
	}
	return nil
}

// governedGenerate 在限流、超时和重试边界内执行分类等非流式调用。
func governedGenerate(
	ctx context.Context,
	provider chatProvider,
	config governanceConfig,
	call modelCall,
) (modelResult, stageMetrics, error) {
	if ctx == nil {
		return modelResult{}, stageMetrics{}, fmt.Errorf("Generate Context 不能为空")
	}
	if provider == nil {
		return modelResult{}, stageMetrics{}, fmt.Errorf("Generate Provider 不能为空")
	}
	if err := validateGovernanceConfig(config); err != nil {
		return modelResult{}, stageMetrics{}, err
	}

	//  12：在每次远程尝试前等待限流器，创建独立超时 Context，并记录 attempts/latency/usage。
	startedAt := time.Now()
	var metrics stageMetrics
	var lastErr error

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		if err := config.Limiter.Wait(ctx); err != nil {
			lastErr = fmt.Errorf("等待模型调用限流失败: %w", err)
			break
		}

		metrics.Attempts++
		attemptCtx, cancel := context.WithTimeout(ctx, config.Timeout)
		result, err := provider.Generate(attemptCtx, call)
		cancel()
		if err == nil {
			metrics.Latency = time.Since(startedAt)
			metrics.Success = true
			metrics.Usage = result.Usage
			return result, metrics, nil
		}

		lastErr = err
		if attempt == config.MaxRetries || !isRetryableModelError(ctx, err) {
			break
		}
		if err := waitBackoff(ctx, config.InitialBackoff, attempt, maximumRetryBackoff); err != nil {
			lastErr = fmt.Errorf("等待第 %d 次重试退避失败: %w", attempt+1, err)
			break
		}
	}

	metrics.Latency = time.Since(startedAt)
	return modelResult{}, metrics, lastErr
}

// governedStream 在首个 delta 前管理流创建重试，并把成功流的 Context 生命周期交给包装器。
func governedStream(
	ctx context.Context,
	provider chatProvider,
	config governanceConfig,
	call modelCall,
) (messageStream, stageMetrics, error) {
	if ctx == nil {
		return nil, stageMetrics{}, fmt.Errorf("Stream Context 不能为空")
	}
	if provider == nil {
		return nil, stageMetrics{}, fmt.Errorf("Stream Provider 不能为空")
	}
	if err := validateGovernanceConfig(config); err != nil {
		return nil, stageMetrics{}, err
	}

	//  13：仅在 Stream 创建失败且未发送 delta 时按白名单重试；成功后不得提前 cancel 流 Context。
	var retries int
	now := time.Now()
	var lastErr error
	var metrics stageMetrics
	for i := 0; i <= config.MaxRetries; i++ {
		if err := config.Limiter.Wait(ctx); err != nil {
			lastErr = fmt.Errorf("等待模型调用限流失败: %w", err)
			break
		}

		retries++
		metrics.Attempts = retries
		timeout, cancelFunc := context.WithTimeout(ctx, config.Timeout)
		stream, err := provider.Stream(timeout, call)
		if err == nil {
			if stream == nil {
				cancelFunc()
				metrics.Latency = time.Since(now)
				return nil, metrics, fmt.Errorf("模型流不能为空")
			}
			metrics.Success = true
			metrics.Latency = time.Since(now)
			return &cancelOnCloseStream{stream: stream, cancel: cancelFunc}, metrics, nil
		}
		cancelFunc()
		lastErr = err
		if !isRetryableModelError(ctx, err) {
			break
		}
		err = waitBackoff(ctx, config.InitialBackoff, retries, maximumRetryBackoff)
		if err != nil {
			lastErr = fmt.Errorf("等待第 %d 次重试退避失败: %w", retries, err)
			break
		}
	}
	metrics.Success = false
	metrics.Latency = time.Since(now)
	return nil, metrics, lastErr
}

// isRetryableModelError 判断错误是否属于允许有限重试的模型服务或临时网络失败。
func isRetryableModelError(parentCtx context.Context, err error) bool {
	if parentCtx == nil || err == nil || parentCtx.Err() != nil {
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
	var temporaryError interface{ Temporary() bool }
	if errors.As(err, &temporaryError) {
		return temporaryError.Temporary()
	}

	//  14：仅允许临时网络错误和 429/500/502/503/504，并实现可取消、有上限的指数退避。
	return false
}

// waitBackoff 执行有上限的指数退避，并在等待期间响应 Context 取消。
func waitBackoff(ctx context.Context, initial time.Duration, retryIndex int, maximumRetryBackoff time.Duration) error {
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
	for i := 0; i < retryIndex && delay < maximumRetryBackoff; i++ {
		if delay > maximumRetryBackoff/2 {
			delay = maximumRetryBackoff
			break
		}
		delay *= 2
	}
	if delay > maximumRetryBackoff {
		delay = maximumRetryBackoff
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
