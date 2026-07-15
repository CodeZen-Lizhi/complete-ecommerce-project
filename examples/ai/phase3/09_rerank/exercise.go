package main

import (
	"context"
	"errors"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type rerankCandidate struct {
	DocumentID string
	Content    string
	BaseScore  float64
}

type rerankResult struct {
	DocumentID string
	Score      float64
}

type reranker interface {
	// Rerank 根据查询重排有限候选集合。
	Rerank(ctx context.Context, query string, candidates []rerankCandidate) ([]rerankResult, error)
}

// validateRerankResults 校验返回 ID、数量、重复项和分数范围。
func validateRerankResults(candidates []rerankCandidate, results []rerankResult) error {
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“候选 Rerank”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：限制候选数量、单段长度和总输入大小。
	// TODO 2：通过 reranker 构造 Query-Document 对并执行重排。
	// TODO 3：实现 validateRerankResults，拒绝未知 ID、重复 ID 和非法分数。
	// TODO 4：按重排分数稳定排序，明确未返回候选的处理策略。
	// TODO 5：记录 MRR/Recall 变化、额外 P95 和 Reranker 失败策略。
	return errExerciseIncomplete
}
