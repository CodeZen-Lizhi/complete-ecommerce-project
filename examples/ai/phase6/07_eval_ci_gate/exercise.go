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

// evaluateGate 返回是否通过和所有阻断原因。
func evaluateGate(input gateInput, thresholds gateThresholds) (bool, []string, error) {
	return false, nil, errExerciseIncomplete
}

// exitCodeForGate 把通过、质量退化和基础设施失败映射为稳定退出码。
func exitCodeForGate(passed bool, evaluationErr error) (int, error) {
	return 0, errExerciseIncomplete
}

// runExercise 按执行顺序组织“CI 质量门禁”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：选择离线、快速、确定性的核心 Case 集合。
	// TODO 2：定义绝对阈值和相对基线退化阈值。
	// TODO 3：实现 evaluateGate，区分质量退化、数据缺失和运行错误。
	// TODO 4：实现 exitCodeForGate，门禁失败必须返回非零退出码。
	// TODO 5：生成机器可读报告、人类摘要和本地复现命令。
	return errExerciseIncomplete
}
