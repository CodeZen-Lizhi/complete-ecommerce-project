package main

import (
	"context"
	"errors"
	"time"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type callMeasurement struct {
	CaseID       string
	Model        string
	StartedAt    time.Time
	FirstTokenAt time.Time
	FinishedAt   time.Time
	PromptTokens int
	OutputTokens int
	Cost         float64
	ErrorType    string
}

// recordCallMeasurement 记录一次调用的时间、Token、成本和错误类型。
func recordCallMeasurement(ctx context.Context) (callMeasurement, error) {
	// TODO 1：采集开始、首 Token、结束时间和用量，保留 Case 与模型标识。
	return callMeasurement{}, errExerciseIncomplete
}

// validateMeasurements 校验时间顺序、用量和样本数量。
func validateMeasurements(values []callMeasurement) error {
	// TODO 2：拒绝时间倒序、负 Token、负成本和空样本。
	return errExerciseIncomplete
}

// percentileDuration 计算排序后样本的指定分位数。
func percentileDuration(samples []time.Duration, percentile float64) (time.Duration, error) {
	// TODO 3：实现分位数算法，并验证 P50/P95 的边界样本。
	return 0, errExerciseIncomplete
}

// summarizeMeasurements 汇总 P50/P95、Token、成本和失败率。
func summarizeMeasurements(values []callMeasurement) error {
	// TODO 4：按模型、Prompt 版本和标签汇总延迟、用量和失败率。
	return errExerciseIncomplete
}

// reportCostQualityTradeoff 同时展示性能、成本和质量指标。
func reportCostQualityTradeoff(ctx context.Context) error {
	// TODO 5：避免只优化成本或延迟而忽略回答质量。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“延迟、Token 与成本指标”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return reportCostQualityTradeoff(ctx)
}
