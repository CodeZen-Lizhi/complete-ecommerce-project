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

// buildEvaluationSnapshots 固定基线和候选配置版本。
func buildEvaluationSnapshots() (evaluationSnapshot, evaluationSnapshot, error) {
	// TODO 1：记录模型、Prompt、切块和检索配置版本。
	return evaluationSnapshot{}, evaluationSnapshot{}, errExerciseIncomplete
}

// validateComparableSnapshots 校验两组运行可直接比较。
func validateComparableSnapshots(baseline evaluationSnapshot, candidate evaluationSnapshot) error {
	// TODO 2：要求相同数据集版本、运行参数和 Case 集合。
	return errExerciseIncomplete
}

// compareSnapshots 使用相同数据集比较基线和候选配置。
func compareSnapshots(baseline evaluationSnapshot, candidate evaluationSnapshot, thresholds regressionThresholds) (regressionReport, error) {
	// TODO 3：计算质量、P95、Token 和成本差值。
	return regressionReport{}, errExerciseIncomplete
}

// applyRegressionThresholds 应用最小样本和允许波动阈值。
func applyRegressionThresholds(report regressionReport, thresholds regressionThresholds) (regressionReport, error) {
	// TODO 4：区分确定退化和样本不足导致的不确定结果。
	return regressionReport{}, errExerciseIncomplete
}

// reportRegressionCases 输出改善、退化和缺失结果 Case。
func reportRegressionCases(ctx context.Context) error {
	// TODO 5：列出具体 Case 和差值证据，不能只给汇总结论。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“回归对比”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return reportRegressionCases(ctx)
}
