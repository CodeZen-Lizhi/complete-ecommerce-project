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

// percentileDuration 计算排序后样本的指定分位数。
func percentileDuration(samples []time.Duration, percentile float64) (time.Duration, error) {
	return 0, errExerciseIncomplete
}

// summarizeMeasurements 汇总 P50/P95、Token、成本和失败率。
func summarizeMeasurements(values []callMeasurement) error {
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“延迟、Token 与成本指标”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：记录开始、首 Token、结束、Token、成本和错误分类。
	// TODO 2：校验时间顺序、Token 非负和样本数量。
	// TODO 3：实现 percentileDuration，并用边界样本验证 P50/P95。
	// TODO 4：实现 summarizeMeasurements，按模型、Prompt 版本和标签汇总。
	// TODO 5：同时展示质量指标，避免只优化成本或延迟。
	return errExerciseIncomplete
}
