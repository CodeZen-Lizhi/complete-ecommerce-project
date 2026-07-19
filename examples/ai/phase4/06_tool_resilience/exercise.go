package main

import (
	"context"
	"errors"
	"time"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type toolFailureKind string

const (
	failureTimeout   toolFailureKind = "timeout"
	failureCanceled  toolFailureKind = "canceled"
	failureTemporary toolFailureKind = "temporary"
	failurePermanent toolFailureKind = "permanent"
)

type resilientTool interface {
	// Invoke 执行幂等只读操作，并返回受限结果。
	Invoke(ctx context.Context, arguments string) (string, error)
}

// invokeWithTimeout 为单次工具调用设置独立超时。
func invokeWithTimeout(ctx context.Context, tool resilientTool, arguments string, timeout time.Duration) (string, error) {
	// TODO 1：派生超时 Context，并确保所有返回路径都会调用 cancel。
	return "", errExerciseIncomplete
}

// classifyToolFailure 把错误分类为取消、超时、临时或永久失败。
func classifyToolFailure(err error) (toolFailureKind, error) {
	// TODO 2：取消和永久业务错误不得标记为可重试。
	return "", errExerciseIncomplete
}

// invokeWithRetry 在独立超时和有限重试内执行只读工具。
func invokeWithRetry(ctx context.Context, tool resilientTool, arguments string, timeout time.Duration, maxRetries int) (string, error) {
	// TODO 3：只对幂等只读临时错误做有限退避重试。
	return "", errExerciseIncomplete
}

// limitToolResponse 限制工具返回的字节数和列表条数。
func limitToolResponse(output string, maxBytes int) (string, error) {
	// TODO 4：结果超限时明确失败，不允许静默截断关键字段。
	return "", errExerciseIncomplete
}

// verifyToolResilience 记录并验证重试治理行为。
func verifyToolResilience(ctx context.Context) error {
	// TODO 5：记录耗时、重试次数和失败类型，并检查资源释放。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“工具失败治理”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return verifyToolResilience(ctx)
}
