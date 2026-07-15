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

// classifyToolFailure 把错误分类为取消、超时、临时或永久失败。
func classifyToolFailure(err error) (toolFailureKind, error) {
	return "", errExerciseIncomplete
}

// invokeWithRetry 在独立超时和有限重试内执行只读工具。
func invokeWithRetry(ctx context.Context, tool resilientTool, arguments string, timeout time.Duration, maxRetries int) (string, error) {
	return "", errExerciseIncomplete
}

// runExercise 按执行顺序组织“工具失败治理”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：为单次工具调用设置独立超时，并确保 cancel 及时执行。
	// TODO 2：实现 classifyToolFailure，取消和永久业务错误不得重试。
	// TODO 3：实现 invokeWithRetry，只对幂等只读临时错误做有限退避重试。
	// TODO 4：限制响应字节数和列表条数，超限明确失败。
	// TODO 5：记录耗时、重试次数和失败类型，并用 race 检查资源释放。
	return errExerciseIncomplete
}
