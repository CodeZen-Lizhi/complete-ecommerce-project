package main

import (
	"context"
	"errors"
	"time"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type ragEvalCase struct {
	ID          string
	Query       string
	RelevantIDs map[string]struct{}
}

type ragEvalRun struct {
	CaseID    string
	RankedIDs []string
	Latency   time.Duration
	Cost      float64
}

// loadRAGEvaluationCases 加载固定问题集和相关 Chunk ID。
func loadRAGEvaluationCases(path string) ([]ragEvalCase, error) {
	// TODO 1：读取过滤条件和配置版本，并拒绝重复 Case ID。
	return nil, errExerciseIncomplete
}

// loadRAGEvaluationRuns 读取不同检索配置的原始排名。
func loadRAGEvaluationRuns(path string) ([]ragEvalRun, error) {
	// TODO 2：保留切块、Hybrid、Rerank 配置和每个 Case 的原始结果。
	return nil, errExerciseIncomplete
}

// recallAtK 计算前 K 个结果覆盖相关文档的比例。
func recallAtK(relevant map[string]struct{}, ranked []string, k int) (float64, error) {
	// TODO 3：实现 Recall@K，并用手算小样本验证公式。
	return 0, errExerciseIncomplete
}

// reciprocalRank 返回首个相关结果的倒数排名。
func reciprocalRank(relevant map[string]struct{}, ranked []string) (float64, error) {
	// TODO 4：返回首个相关结果的倒数排名，并处理无命中情况。
	return 0, errExerciseIncomplete
}

// summarizeRAGEvaluation 汇总 Recall@K、MRR、延迟和成本。
func summarizeRAGEvaluation(cases []ragEvalCase, runs []ragEvalRun, k int) error {
	// TODO 5：计算分组 Recall、MRR、P50/P95 和成本。
	return errExerciseIncomplete
}

// reportRAGEvaluationFailures 输出具体失败和缺失结果 Case。
func reportRAGEvaluationFailures(ctx context.Context) error {
	// TODO 6：列出失败 Case，缺失运行结果时明确失败。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“RAG 检索评估”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return reportRAGEvaluationFailures(ctx)
}
