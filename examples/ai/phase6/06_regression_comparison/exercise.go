package main

import (
	"context"
	"errors"
	"time"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type evaluationSnapshot struct {
	ConfigVersion  string
	DatasetVersion string
	Quality        float64
	P95            time.Duration
	Tokens         int
	Cost           float64
	Failures       map[string]string
}

type regressionThresholds struct {
	MaxQualityDrop  float64
	MaxP95Increase  time.Duration
	MaxCostIncrease float64
}

type regressionReport struct {
	Improved  []string
	Regressed []string
	Uncertain []string
}

// compareSnapshots 使用相同数据集比较基线和候选配置。
func compareSnapshots(baseline evaluationSnapshot, candidate evaluationSnapshot, thresholds regressionThresholds) (regressionReport, error) {
	return regressionReport{}, errExerciseIncomplete
}

// runExercise 按执行顺序组织“回归对比”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：固定基线与候选的模型、Prompt、切块和检索配置版本。
	// TODO 2：验证两组使用相同数据集版本和运行参数。
	// TODO 3：实现 compareSnapshots，计算质量、P95、Token 和成本差值。
	// TODO 4：应用最小样本和允许波动阈值，区分退化与不确定。
	// TODO 5：列出具体改善、退化和缺失结果 Case。
	return errExerciseIncomplete
}
