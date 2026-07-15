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

// decideFallback 根据失败类型和请求副作用边界决定是否降级。
func decideFallback(failure dependencyFailure, request fallbackRequest) (fallbackDecision, error) {
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

// runExercise 按执行顺序组织“依赖故障与降级”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：定义可降级和不可降级场景，写操作禁止假成功。
	// TODO 2：注入模型超时、429、5xx 和向量库连接失败。
	// TODO 3：实现 decideFallback，只允许安全只读请求使用缓存或备用模型。
	// TODO 4：使用 recoveryController 限制恢复探测频率，避免重试风暴。
	// TODO 5：验证降级响应标记能力限制，依赖恢复后自动退出。
	return errExerciseIncomplete
}
