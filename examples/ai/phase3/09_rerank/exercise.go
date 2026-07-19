package main

import (
	"context"
	"errors"
)

// Reranker 配置集中放在顶部，练习时直接替换占位值。
const (
	rerankerBaseURL   = "http://localhost:8084/v1"
	rerankerAPIKey    = "replace-with-your-api-key"
	rerankerModelName = "bge-reranker-v2-m3"
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

type rerankerConfig struct {
	BaseURL string
	Model   string
	APIKey  string
}

// validateRerankCandidates 校验远程重排的输入预算。
func validateRerankCandidates(candidates []rerankCandidate) error {
	// TODO 1：限制候选数量、单段长度和总输入大小。
	return errExerciseIncomplete
}

// newRemoteReranker 从配置创建真实远程 Reranker Adapter。
func newRemoteReranker(ctx context.Context, config rerankerConfig) (reranker, error) {
	// TODO 2：构造 Query-Document 对并调用真实 Reranker 服务。
	return nil, errExerciseIncomplete
}

// validateRerankResults 校验返回 ID、数量、重复项和分数范围。
func validateRerankResults(candidates []rerankCandidate, results []rerankResult) error {
	// TODO 3：拒绝未知 ID、重复 ID、数量异常和非法分数。
	return errExerciseIncomplete
}

// stableSortRerankResults 按重排分数稳定排序候选。
func stableSortRerankResults(candidates []rerankCandidate, results []rerankResult) ([]rerankCandidate, error) {
	// TODO 4：稳定排序，并明确未返回候选的保留或丢弃策略。
	return nil, errExerciseIncomplete
}

// summarizeRerankEvaluation 汇总重排前后的质量和延迟。
func summarizeRerankEvaluation(ctx context.Context, config rerankerConfig) error {
	// TODO 5：记录 MRR/Recall 变化、额外 P95 和 Reranker 失败策略。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“候选 Rerank”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return summarizeRerankEvaluation(ctx, rerankerConfig{
		BaseURL: rerankerBaseURL,
		Model:   rerankerModelName,
		APIKey:  rerankerAPIKey,
	})
}
