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

// recallAtK 计算前 K 个结果覆盖相关文档的比例。
func recallAtK(relevant map[string]struct{}, ranked []string, k int) (float64, error) {
	return 0, errExerciseIncomplete
}

// reciprocalRank 返回首个相关结果的倒数排名。
func reciprocalRank(relevant map[string]struct{}, ranked []string) (float64, error) {
	return 0, errExerciseIncomplete
}

// summarizeRAGEvaluation 汇总 Recall@K、MRR、延迟和成本。
func summarizeRAGEvaluation(cases []ragEvalCase, runs []ragEvalRun, k int) error {
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“RAG 检索评估”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：加载固定问题集、相关 Chunk ID、过滤条件和配置版本。
	// TODO 2：保存不同切块、Hybrid、Rerank 配置的原始排名。
	// TODO 3：实现 recallAtK 和 reciprocalRank，并用小样本验证公式。
	// TODO 4：实现 summarizeRAGEvaluation，计算分组指标、P50/P95 和成本。
	// TODO 5：输出具体失败 case，缺失运行结果时明确失败。
	return errExerciseIncomplete
}
