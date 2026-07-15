package main

import (
	"context"
	"errors"
	"time"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type metricLabels map[string]string

type metricRecorder interface {
	// ObserveDuration 记录耗时直方图。
	ObserveDuration(name string, value time.Duration, labels metricLabels)
	// AddCounter 累加请求、错误、Token 或成本计数。
	AddCounter(name string, value float64, labels metricLabels)
	// SetGauge 设置队列长度或并发量。
	SetGauge(name string, value float64, labels metricLabels)
}

type alertRule struct {
	Metric    string
	Threshold float64
	Window    time.Duration
}

// validateMetricLabels 拒绝 user ID、session ID、Prompt 等高基数标签。
func validateMetricLabels(labels metricLabels) error {
	return errExerciseIncomplete
}

// evaluateAlert 判断窗口指标是否触发或恢复告警。
func evaluateAlert(rule alertRule, values []float64) (bool, error) {
	return false, errExerciseIncomplete
}

// runExercise 按执行顺序组织“Metrics 与告警”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：定义请求、成功率、P50/P95、Token、成本和队列指标。
	// TODO 2：实现 validateMetricLabels，拒绝高基数和敏感标签。
	// TODO 3：使用 metricRecorder 记录模型、检索和 Agent 指标。
	// TODO 4：实现 evaluateAlert，为错误率、延迟、积压和预算设置窗口阈值。
	// TODO 5：用本地 fixture 验证告警触发、持续和恢复。
	return errExerciseIncomplete
}
