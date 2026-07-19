package main

import (
	"context"
	"errors"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type retrievalCase struct {
	ID          string
	RelevantIDs map[string]struct{}
}

type retrievalRun struct {
	CaseID    string
	RankedIDs []string
}

// loadRetrievalEvaluation 加载 Golden Dataset 和原始排名。
func loadRetrievalEvaluation() ([]retrievalCase, []retrievalRun, error) {
	// TODO 1：对齐 Case ID，并保留每个查询的完整原始排名。
	return nil, nil, errExerciseIncomplete
}

// validateRetrievalInputs 校验指标输入边界。
func validateRetrievalInputs(cases []retrievalCase, runs []retrievalRun, k int) error {
	// TODO 2：拒绝非法 K、空相关集合、重复结果和缺失 Case。
	return errExerciseIncomplete
}

// recallAtK 计算前 K 个结果中相关文档的覆盖比例。
func recallAtK(relevant map[string]struct{}, ranked []string, k int) (float64, error) {
	// TODO 3：实现 Recall@K，并用手算样本验证无命中和截断边界。
	return 0, errExerciseIncomplete
}

// reciprocalRank 计算首个相关结果的倒数排名。
func reciprocalRank(relevant map[string]struct{}, ranked []string) (float64, error) {
	// TODO 4：实现首个相关结果的倒数排名，无命中时返回零。
	return 0, errExerciseIncomplete
}

// summarizeRetrieval 汇总总体与标签分组的 Recall@K、MRR 和 Hit Rate。
func summarizeRetrieval(cases []retrievalCase, runs []retrievalRun, k int) error {
	// TODO 5：计算总体和分组 Hit Rate、Recall 与 MRR。
	return errExerciseIncomplete
}

// reportRetrievalFailures 输出失败 Case 和排名证据。
func reportRetrievalFailures(ctx context.Context) error {
	// TODO 6：列出具体失败 Case，禁止只输出平均数。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“检索指标”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return reportRetrievalFailures(ctx)
}
