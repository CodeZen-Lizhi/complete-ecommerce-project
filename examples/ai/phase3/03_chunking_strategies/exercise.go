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

// validateChunkingConfig 校验块大小和重叠边界。
func validateChunkingConfig(config chunkingConfig) error {
	// TODO 1：要求 MaxRunes 为正、Overlap 非负且小于块大小。
	return errExerciseIncomplete
}

// splitByLength 实现递归长度切块。
func splitByLength(blocks []sourceBlock, config chunkingConfig) ([]textChunk, error) {
	// TODO 2：保持原文顺序，并正确处理 overlap 边界。
	return nil, errExerciseIncomplete
}

// splitByStructure 按标题、段落和列表执行结构化切块。
func splitByStructure(blocks []sourceBlock, config chunkingConfig) ([]textChunk, error) {
	// TODO 3：保留结构类型和标题上下文，避免跨语义边界硬切。
	return nil, errExerciseIncomplete
}

// buildParentChildChunks 生成父块和可索引子块，并保持稳定关联 ID。
func buildParentChildChunks(blocks []sourceBlock, config chunkingConfig) ([]textChunk, error) {
	// TODO 4：生成稳定父子 ID，并为表格和代码块选择独立策略。
	return nil, errExerciseIncomplete
}

// compareChunkingStrategies 输出不同策略的块分布和边界。
func compareChunkingStrategies(ctx context.Context) error {
	// TODO 5：输出块数量、长度分布、边界内容和父子关系。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“切块策略对比”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return compareChunkingStrategies(ctx)
}
