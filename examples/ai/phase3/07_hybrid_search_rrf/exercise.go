package main

import (
	"context"
	"errors"
)

// 模型与检索组件配置集中放在顶部，练习时直接替换占位值。
const (
	embeddingBaseURL   = "http://localhost:8084/v1"
	embeddingAPIKey    = "replace-with-your-api-key"
	embeddingModelName = "text-embedding-3-small"
	qdrantBaseURL      = "http://localhost:6333"
	qdrantAPIKey       = "replace-with-your-qdrant-api-key"
	qdrantCollection   = "ai_phase3_hybrid"
	bm25BaseURL        = "http://localhost:8085"
	bm25APIKey         = "replace-with-your-bm25-api-key"
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
	EmbeddingBaseURL string
	EmbeddingAPIKey  string
	EmbeddingModel   string
	QdrantBaseURL    string
	QdrantAPIKey     string
	BM25BaseURL      string
	BM25APIKey       string
	Collection       string
}

// validateHybridContract 校验两路检索共享的文档与过滤契约。
func validateHybridContract(topK int, documentIDs []string) error {
	// TODO 1：确保 Dense 与 BM25 使用相同文档 ID、权限过滤和 Top-K。
	return errExerciseIncomplete
}

// newHybridRetrievers 连接真实 Qdrant Dense 后端和 BM25 Sparse 后端。
func newHybridRetrievers(ctx context.Context, config hybridBackendConfig) (denseRetriever, sparseRetriever, error) {
	// TODO 2：连接真实 Qdrant 与 BM25 后端，并保留两路原始排名和耗时。
	return nil, nil, errExerciseIncomplete
}

// reciprocalRankFusion 使用 1/(k+rank) 融合多路排名并按文档去重。
func reciprocalRankFusion(rankings [][]rankedHit, k int, topK int) ([]rankedHit, error) {
	// TODO 3：校验 rank、k、topK，并按文档 ID 累加 RRF 分数。
	return nil, errExerciseIncomplete
}

// stableSortFusedHits 对融合分数执行稳定排序。
func stableSortFusedHits(hits []rankedHit) ([]rankedHit, error) {
	// TODO 4：稳定处理并列分数和只在单路出现的文档。
	return nil, errExerciseIncomplete
}

// compareHybridResults 对比三种检索结果的质量与延迟。
func compareHybridResults(ctx context.Context, config hybridBackendConfig) error {
	// TODO 5：对比 Dense、BM25 和融合结果的命中率与延迟。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“Dense + BM25 + RRF”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return compareHybridResults(ctx, hybridBackendConfig{
		EmbeddingBaseURL: embeddingBaseURL,
		EmbeddingAPIKey:  embeddingAPIKey,
		EmbeddingModel:   embeddingModelName,
		QdrantBaseURL:    qdrantBaseURL,
		QdrantAPIKey:     qdrantAPIKey,
		BM25BaseURL:      bm25BaseURL,
		BM25APIKey:       bm25APIKey,
		Collection:       qdrantCollection,
	})
}
