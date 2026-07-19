package main

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

// 外部组件配置集中放在顶部，练习时直接替换占位值。
const (
	embeddingBaseURL   = "http://localhost:8084/v1"
	embeddingAPIKey    = "replace-with-your-api-key"
	embeddingModelName = "text-embedding-3-small"
	qdrantBaseURL      = "http://localhost:6333"
	qdrantAPIKey       = "replace-with-your-qdrant-api-key"
	qdrantCollection   = "ai_phase3_documents"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type indexedDocument struct {
	ID       string
	Content  string
	Metadata map[string]string
}

type indexer interface {
	// Index 批量写入文档；重复 ID 的处理策略必须明确。
	Index(ctx context.Context, documents []indexedDocument) error
	// Delete 按文档 ID 删除索引内容。
	Delete(ctx context.Context, ids []string) error
}

type retriever interface {
	// Retrieve 根据查询和 Top-K 返回按相关度排序的文档。
	Retrieve(ctx context.Context, query string, topK int) ([]indexedDocument, error)
}

type qdrantConfig struct {
	BaseURL          string
	Collection       string
	APIKey           string
	EmbeddingBaseURL string
	EmbeddingAPIKey  string
	EmbeddingModel   string
}

// buildIndexDocuments 定义稳定 ID、向量维度和 Metadata 契约。
func buildIndexDocuments() ([]indexedDocument, error) {
	// TODO 1：准备带稳定 Document ID、Chunk ID 和权限 Metadata 的文档批次。
	return nil, errExerciseIncomplete
}

// newQdrantHTTPClient 校验真实 Qdrant 地址和集合配置，并返回用于 REST API 调用的 Client。
func newQdrantHTTPClient(config qdrantConfig) (*http.Client, error) {
	if strings.TrimSpace(config.BaseURL) == "" || strings.TrimSpace(config.Collection) == "" {
		return nil, errors.New("Qdrant Base URL 和 Collection 不能为空")
	}
	// TODO 2：创建真实 Qdrant Client 和集合，再接入 Eino Embedder、Indexer 与 Retriever。
	return nil, errExerciseIncomplete
}

// validateIndexBatch 校验 ID、内容、Metadata 和批量上限。
func validateIndexBatch(documents []indexedDocument) error {
	// TODO 3：拒绝空 ID、空内容、缺失 Metadata、重复 ID 和超大批次。
	return errExerciseIncomplete
}

// retrieveAndValidate 执行检索并校验 Top-K、空结果和返回顺序。
func retrieveAndValidate(ctx context.Context, retriever retriever, query string, topK int) ([]indexedDocument, error) {
	// TODO 4：调用 Retrieve，并验证返回数量、顺序和 Context 取消。
	return nil, errExerciseIncomplete
}

// verifyIndexLifecycle 验证真实 Qdrant 的写入、查询和删除闭环。
func verifyIndexLifecycle(ctx context.Context, config qdrantConfig) error {
	// TODO 5：输出命中文档、分数和 Metadata，并验证 Delete 后不再召回。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“向量索引与检索”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return verifyIndexLifecycle(ctx, qdrantConfig{
		BaseURL:          qdrantBaseURL,
		Collection:       qdrantCollection,
		APIKey:           qdrantAPIKey,
		EmbeddingBaseURL: embeddingBaseURL,
		EmbeddingAPIKey:  embeddingAPIKey,
		EmbeddingModel:   embeddingModelName,
	})
}
