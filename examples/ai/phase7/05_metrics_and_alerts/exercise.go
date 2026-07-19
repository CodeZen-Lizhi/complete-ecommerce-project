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

// defineAIMetrics 定义请求、质量、用量、成本和队列指标。
func defineAIMetrics() ([]string, error) {
	// TODO 1：为每个指标明确类型、单位和允许标签。
	return nil, errExerciseIncomplete
}

// validateMetricLabels 拒绝 user ID、session ID、Prompt 等高基数标签。
func validateMetricLabels(labels metricLabels) error {
	// TODO 2：拒绝高基数、敏感和未注册标签。
	return errExerciseIncomplete
}

// recordPipelineMetrics 使用 Recorder 写入模型、检索和 Agent 指标。
func recordPipelineMetrics(ctx context.Context, recorder metricRecorder) error {
	// TODO 3：记录耗时、请求数、错误、Token、成本和队列长度。
	return errExerciseIncomplete
}

// evaluateAlert 判断窗口指标是否触发或恢复告警。
func evaluateAlert(rule alertRule, values []float64) (bool, error) {
	// TODO 4：为错误率、延迟、积压和预算设置窗口阈值与恢复条件。
	return false, errExerciseIncomplete
}

// verifyAlertLifecycle 验证告警触发、持续和恢复。
func verifyAlertLifecycle(ctx context.Context) error {
	// TODO 5：使用本地样本覆盖阈值边界和空窗口。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“Metrics 与告警”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return verifyAlertLifecycle(ctx)
}
