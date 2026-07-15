package main

import (
	"context"
	"errors"
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

type childRetriever interface {
	// RetrieveChildren 返回按分数排序的子块命中。
	RetrieveChildren(ctx context.Context, query string, topK int) ([]childHit, error)
}

type parentStore interface {
	// LoadParents 批量读取父块；返回值必须能按 ID 定位。
	LoadParents(ctx context.Context, parentIDs []string) (map[string]parentDocument, error)
}

// expandParents 合并同一父块的多个子块命中，并保留最佳分数。
func expandParents(hits []childHit, parents map[string]parentDocument) ([]parentDocument, error) {
	return nil, errExerciseIncomplete
}

// runExercise 按执行顺序组织“Parent-Child Retrieval”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：导入时生成稳定 parent ID 和 child ID，并让子块继承权限 Metadata。
	// TODO 2：只索引子块，通过 childRetriever 执行 Top-K 召回。
	// TODO 3：去重 parent ID 后使用 parentStore 批量加载父块，禁止循环逐条查询。
	// TODO 4：实现 expandParents，缺失父块明确失败并保留最佳子块分数。
	// TODO 5：验证一父多子、权限继承、空结果和父块缺失场景。
	return errExerciseIncomplete
}
