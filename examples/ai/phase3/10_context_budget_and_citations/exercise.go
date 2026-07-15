package main

import (
	"context"
	"errors"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type contextCandidate struct {
	DocumentID string
	ParentID   string
	FileName   string
	Page       int
	Content    string
	Score      float64
}

type selectedContext struct {
	CitationID string
	Candidate  contextCandidate
	Tokens     int
}

type tokenCounter interface {
	// Count 返回文本在目标模型下的 Token 数量。
	Count(ctx context.Context, text string) (int, error)
}

// selectContext 在预算内去重、扩展父块并分配稳定引用 ID。
func selectContext(ctx context.Context, counter tokenCounter, candidates []contextCandidate, budget int) ([]selectedContext, error) {
	return nil, errExerciseIncomplete
}

// runExercise 按执行顺序组织“上下文预算与引用”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：按文档、页码和位置去重重叠候选。
	// TODO 2：批量扩展父块，避免重复加载同一父块。
	// TODO 3：通过 tokenCounter 计算每段成本并校验正预算。
	// TODO 4：实现 selectContext，按相关度在预算内选择并分配稳定 Citation ID。
	// TODO 5：验证回答引用只指向实际选中的文件名、页码和证据。
	return errExerciseIncomplete
}
