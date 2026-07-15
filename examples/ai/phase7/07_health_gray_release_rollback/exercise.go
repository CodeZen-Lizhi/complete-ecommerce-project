package main

import (
	"context"
	"errors"
	"time"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type healthReport struct {
	Live     bool
	Ready    bool
	Degraded bool
	Reason   string
}

type releaseMetrics struct {
	Version      string
	RequestCount int
	ErrorRate    float64
	P95          time.Duration
	Quality      float64
}

type rolloutDecision string

const (
	rolloutHold     rolloutDecision = "hold"
	rolloutExpand   rolloutDecision = "expand"
	rolloutRollback rolloutDecision = "rollback"
)

// evaluateRollout 根据错误率、P95、质量和最小样本决定扩量或回滚。
func evaluateRollout(baseline releaseMetrics, candidate releaseMetrics) (rolloutDecision, error) {
	return rolloutHold, errExerciseIncomplete
}

// validateHealth 区分 liveness、readiness 和依赖降级。
func validateHealth(report healthReport) error {
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“健康检查、灰度与回滚”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：定义 liveness、readiness 和 degraded 的独立语义。
	// TODO 2：实现 validateHealth，依赖降级不能伪装为完全健康。
	// TODO 3：定义版本、配置和模型路由发布标识。
	// TODO 4：实现 evaluateRollout，最小样本不足时保持，不盲目扩量。
	// TODO 5：验证回滚后会话、任务和索引版本仍兼容。
	return errExerciseIncomplete
}
