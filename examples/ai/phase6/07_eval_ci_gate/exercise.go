package main

import (
	"context"
	"errors"
	"time"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type gateThresholds struct {
	MinQuality     float64
	MaxQualityDrop float64
	MaxP95         time.Duration
	MaxFailureRate float64
}

type gateInput struct {
	BaselineVersion  string
	CandidateVersion string
	Quality          float64
	QualityDrop      float64
	P95              time.Duration
	FailureRate      float64
	MissingCases     int
}

// loadCoreGateCases 选择离线、快速且确定的核心样本。
func loadCoreGateCases() ([]string, error) {
	// TODO 1：固定 Case ID、数据版本和本地复现输入。
	return nil, errExerciseIncomplete
}

// validateGateThresholds 校验绝对阈值和相对退化阈值。
func validateGateThresholds(thresholds gateThresholds) error {
	// TODO 2：拒绝负阈值、互相矛盾的边界和缺失基线。
	return errExerciseIncomplete
}

// evaluateGate 返回是否通过和所有阻断原因。
func evaluateGate(input gateInput, thresholds gateThresholds) (bool, []string, error) {
	// TODO 3：区分质量退化、数据缺失和运行错误，并收集所有阻断原因。
	return false, nil, errExerciseIncomplete
}

// exitCodeForGate 把通过、质量退化和基础设施失败映射为稳定退出码。
func exitCodeForGate(passed bool, evaluationErr error) (int, error) {
	// TODO 4：为通过、质量失败和基础设施失败返回稳定退出码。
	return 0, errExerciseIncomplete
}

// writeGateReport 生成机器报告、人类摘要和复现命令。
func writeGateReport(ctx context.Context) error {
	// TODO 5：门禁失败必须有非零退出码和可执行的本地复现步骤。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“CI 质量门禁”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return writeGateReport(ctx)
}
