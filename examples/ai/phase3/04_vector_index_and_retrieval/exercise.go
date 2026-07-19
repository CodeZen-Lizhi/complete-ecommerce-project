package main

import (
	"context"
	"errors"
	"net/http"
	"strings"
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
	BaseURL    string
	Collection string
	APIKey     string
}

// newQdrantHTTPClient 校验真实 Qdrant 地址和集合配置，并返回用于 REST API 调用的 Client。
func newQdrantHTTPClient(config qdrantConfig) (*http.Client, error) {
	if strings.TrimSpace(config.BaseURL) == "" || strings.TrimSpace(config.Collection) == "" {
		return nil, errors.New("Qdrant Base URL 和 Collection 不能为空")
	}
	return nil, errExerciseIncomplete
}

// validateIndexBatch 校验 ID、内容、Metadata 和批量上限。
func validateIndexBatch(documents []indexedDocument) error {
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“向量索引与检索”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：定义稳定文档 ID、Chunk ID、向量维度和 Metadata 契约。
	// TODO 2：从 VECTOR_BASE_URL 创建真实 Qdrant Client 和集合，再接入 Eino Embedder、Indexer 与 Retriever。
	// TODO 3：实现 validateIndexBatch 后批量写入，区分重复 ID、部分失败和完整失败。
	// TODO 4：对查询执行 Retrieve，校验 Top-K、空结果和返回顺序。
	// TODO 5：对真实 Qdrant 输出命中文档、分数和 Metadata，并验证 Delete 后不再召回。
	return errExerciseIncomplete
}
