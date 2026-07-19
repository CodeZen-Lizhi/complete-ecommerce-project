package main

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"
)

// 发布控制配置集中放在顶部，练习时直接替换占位值。
const (
	releaseControlBaseURL = "http://localhost:8090"
	releaseControlAPIKey  = "replace-with-your-release-control-api-key"
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

type releaseControlConfig struct {
	BaseURL string
	APIKey  string
}

// buildHealthReport 分别采集存活、就绪和依赖降级状态。
func buildHealthReport(ctx context.Context) (healthReport, error) {
	// TODO 1：定义 liveness、readiness 和 degraded 的独立语义。
	return healthReport{}, errExerciseIncomplete
}

// validateHealth 区分 liveness、readiness 和依赖降级。
func validateHealth(report healthReport) error {
	// TODO 2：依赖降级不能伪装为完全健康，并保留可诊断原因。
	return errExerciseIncomplete
}

// newReleaseControlClient 校验真实发布控制端点，并返回执行流量调整和回滚请求的 HTTP Client。
func newReleaseControlClient(config releaseControlConfig) (*http.Client, error) {
	if strings.TrimSpace(config.BaseURL) == "" {
		return nil, errors.New("发布控制 Base URL 不能为空")
	}
	// TODO 3：创建 Client，并读取版本、配置和模型路由标识。
	return nil, errExerciseIncomplete
}

// evaluateRollout 根据错误率、P95、质量和最小样本决定扩量或回滚。
func evaluateRollout(baseline releaseMetrics, candidate releaseMetrics) (rolloutDecision, error) {
	// TODO 4：样本不足时保持；达到阈值时返回 Expand 或 Rollback 决策。
	return rolloutHold, errExerciseIncomplete
}

// runRolloutDrill 执行一次真实扩量和回滚演练。
func runRolloutDrill(ctx context.Context, config releaseControlConfig) error {
	// TODO 5：调用流量调整和回滚端点，并验证会话、任务和索引版本兼容。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“健康检查、灰度与回滚”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return runRolloutDrill(ctx, releaseControlConfig{
		BaseURL: releaseControlBaseURL,
		APIKey:  releaseControlAPIKey,
	})
}
