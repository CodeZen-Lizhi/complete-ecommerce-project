package main

import (
	"context"
	"errors"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type dependencyFailure string

const (
	failureModelTimeout   dependencyFailure = "model_timeout"
	failureModelRateLimit dependencyFailure = "model_rate_limit"
	failureVectorStore    dependencyFailure = "vector_store"
)

type fallbackRequest struct {
	ReadOnly            bool
	CanUseCache         bool
	CanUseFallbackModel bool
}

type fallbackDecision struct {
	Allowed bool
	Mode    string
	Reason  string
}

// defineFallbackBoundaries 列出可降级和不可降级请求。
func defineFallbackBoundaries() ([]fallbackRequest, error) {
	// TODO 1：写操作和无法标记能力限制的请求禁止假成功。
	return nil, errExerciseIncomplete
}

// injectDependencyFailure 注入模型和向量库故障。
func injectDependencyFailure(ctx context.Context, failure dependencyFailure) error {
	// TODO 2：覆盖模型超时、429、5xx 和向量库连接失败。
	return errExerciseIncomplete
}

// decideFallback 根据失败类型和请求副作用边界决定是否降级。
func decideFallback(failure dependencyFailure, request fallbackRequest) (fallbackDecision, error) {
	// TODO 3：只允许安全只读请求使用缓存或备用模型。
	return fallbackDecision{}, errExerciseIncomplete
}

// recoveryController 控制降级进入、探测恢复和退出，避免重试风暴。
type recoveryController interface {
	// Enter 记录降级原因和开始时间。
	Enter(ctx context.Context, failure dependencyFailure) error
	// Probe 执行有频率上限的恢复探测。
	Probe(ctx context.Context) (bool, error)
	// Exit 退出降级并清理临时状态。
	Exit(ctx context.Context) error
}

// controlDependencyRecovery 限制恢复探测并安全退出降级。
func controlDependencyRecovery(ctx context.Context, controller recoveryController, failure dependencyFailure) error {
	// TODO 4：限制 Probe 频率，避免并发探测和重试风暴。
	return errExerciseIncomplete
}

// verifyFallbackLifecycle 验证降级标记和自动恢复。
func verifyFallbackLifecycle(ctx context.Context) error {
	// TODO 5：响应必须标记能力限制，依赖恢复后自动退出降级。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“依赖故障与降级”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return verifyFallbackLifecycle(ctx)
}
