package main

import (
	"context"
	"errors"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type rankedHit struct {
	DocumentID string
	Rank       int
	Score      float64
	Source     string
}

type denseRetriever interface {
	// DenseSearch 返回向量检索排名。
	DenseSearch(ctx context.Context, query string, topK int) ([]rankedHit, error)
}

type sparseRetriever interface {
	// SparseSearch 返回 BM25 检索排名。
	SparseSearch(ctx context.Context, query string, topK int) ([]rankedHit, error)
}

type hybridBackendConfig struct {
	QdrantBaseURL string
	BM25BaseURL   string
	Collection    string
}

// newHybridRetrievers 连接真实 Qdrant Dense 后端和 BM25 Sparse 后端。
func newHybridRetrievers(ctx context.Context, config hybridBackendConfig) (denseRetriever, sparseRetriever, error) {
	return nil, nil, errExerciseIncomplete
}

// reciprocalRankFusion 使用 1/(k+rank) 融合多路排名并按文档去重。
func reciprocalRankFusion(rankings [][]rankedHit, k int, topK int) ([]rankedHit, error) {
	return nil, errExerciseIncomplete
}

// runExercise 按执行顺序组织“Dense + BM25 + RRF”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：让 Dense 与 BM25 使用相同文档 ID、权限过滤和 Top-K 契约。
	// TODO 2：实现 newHybridRetrievers，调用真实 Qdrant Dense 与 BM25 Sparse 检索并保存原始排名与耗时。
	// TODO 3：实现 reciprocalRankFusion，校验 rank、k、topK 并按文档 ID 累加分数。
	// TODO 4：稳定处理并列分数和只在单路出现的文档。
	// TODO 5：对比 Dense、BM25 和融合结果的命中率与延迟。
	return errExerciseIncomplete
}
