package main

import (
	"context"
	"fmt"
	"time"

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

	// TODO 12：在每次远程尝试前等待限流器，创建独立超时 Context，并记录 attempts/latency/usage。
	return modelResult{}, stageMetrics{}, errExerciseIncomplete
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

	// TODO 13：仅在 Stream 创建失败且未发送 delta 时按白名单重试；成功后不得提前 cancel 流 Context。
	return nil, stageMetrics{}, errExerciseIncomplete
}

// isRetryableModelError 判断错误是否属于允许有限重试的模型服务或临时网络失败。
func isRetryableModelError(parentCtx context.Context, err error) bool {
	if parentCtx == nil || err == nil || parentCtx.Err() != nil {
		return false
	}

	// TODO 14：仅允许临时网络错误和 429/500/502/503/504，并实现可取消、有上限的指数退避。
	return false
}
