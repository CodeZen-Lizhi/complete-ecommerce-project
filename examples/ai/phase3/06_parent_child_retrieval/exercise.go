package main

import (
	"context"
	"errors"
)

// Qdrant 配置集中放在顶部，练习时直接替换占位值。
const (
	qdrantBaseURL    = "http://localhost:6333"
	qdrantAPIKey     = "replace-with-your-qdrant-api-key"
	qdrantCollection = "ai_phase3_parent_child"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type childHit struct {
	ChildID  string
	ParentID string
	Score    float64
}

type parentDocument struct {
	ID       string
	Content  string
	Metadata map[string]string
}

type parentChildBackendConfig struct {
	BaseURL    string
	APIKey     string
	Collection string
}

type childRetriever interface {
	// RetrieveChildren 返回按分数排序的子块命中。
	RetrieveChildren(ctx context.Context, query string, topK int) ([]childHit, error)
}

type parentStore interface {
	// LoadParents 批量读取父块；返回值必须能按 ID 定位。
	LoadParents(ctx context.Context, parentIDs []string) (map[string]parentDocument, error)
}

// buildParentChildIndex 生成稳定父子 ID 并继承权限 Metadata。
func buildParentChildIndex(ctx context.Context) error {
	// TODO 1：导入父块和子块，只把子块写入检索索引。
	return errExerciseIncomplete
}

// retrieveChildHits 执行子块 Top-K 召回。
func retrieveChildHits(ctx context.Context, retriever childRetriever, query string, topK int) ([]childHit, error) {
	// TODO 2：调用 RetrieveChildren，并校验返回顺序和权限范围。
	return nil, errExerciseIncomplete
}

// loadUniqueParents 去重 Parent ID 后批量加载父块。
func loadUniqueParents(ctx context.Context, store parentStore, hits []childHit) (map[string]parentDocument, error) {
	// TODO 3：禁止循环逐条查询，并拒绝跨权限父块。
	return nil, errExerciseIncomplete
}

// expandParents 合并同一父块的多个子块命中，并保留最佳分数。
func expandParents(hits []childHit, parents map[string]parentDocument) ([]parentDocument, error) {
	// TODO 4：缺失父块明确失败，并为每个父块保留最佳子块分数。
	return nil, errExerciseIncomplete
}

// verifyParentChildRetrieval 验证父子检索的关键边界。
func verifyParentChildRetrieval(ctx context.Context, backend parentChildBackendConfig) error {
	// TODO 5：覆盖一父多子、权限继承、空结果和父块缺失场景。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“Parent-Child Retrieval”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return verifyParentChildRetrieval(ctx, parentChildBackendConfig{
		BaseURL:    qdrantBaseURL,
		APIKey:     qdrantAPIKey,
		Collection: qdrantCollection,
	})
}
