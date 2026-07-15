package main

import (
	"context"
	"errors"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type chunkingConfig struct {
	MaxRunes int
	Overlap  int
}

type sourceBlock struct {
	Kind    string
	Heading string
	Content string
}

type textChunk struct {
	ID       string
	ParentID string
	Kind     string
	Content  string
}

type chunker interface {
	// Split 根据配置把结构块切成有序 Chunk。
	Split(ctx context.Context, blocks []sourceBlock, config chunkingConfig) ([]textChunk, error)
}

// buildParentChildChunks 生成父块和可索引子块，并保持稳定关联 ID。
func buildParentChildChunks(blocks []sourceBlock, config chunkingConfig) ([]textChunk, error) {
	return nil, errExerciseIncomplete
}

// runExercise 按执行顺序组织“切块策略对比”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：校验 MaxRunes 为正、Overlap 非负且小于块大小。
	// TODO 2：实现递归长度切块，保证顺序和 overlap 边界。
	// TODO 3：实现按标题、段落和列表的结构化切块。
	// TODO 4：实现 buildParentChildChunks，并为表格和代码块选择独立策略。
	// TODO 5：输出块数量、长度分布、边界内容和父子关系用于对比。
	return errExerciseIncomplete
}
