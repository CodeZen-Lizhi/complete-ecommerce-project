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

// deduplicateCandidates 按文档、页码和位置去重重叠候选。
func deduplicateCandidates(candidates []contextCandidate) ([]contextCandidate, error) {
	// TODO 1：稳定去重，并保留得分最高且 Metadata 完整的候选。
	return nil, errExerciseIncomplete
}

// expandCandidateParents 批量扩展父块。
func expandCandidateParents(ctx context.Context, candidates []contextCandidate) ([]contextCandidate, error) {
	// TODO 2：避免重复加载同一父块，并保持权限边界。
	return nil, errExerciseIncomplete
}

// countCandidateTokens 计算每段候选的 Token 成本。
func countCandidateTokens(ctx context.Context, counter tokenCounter, candidates []contextCandidate) ([]int, error) {
	// TODO 3：调用 tokenCounter，并拒绝非正预算和非法 Token 数。
	return nil, errExerciseIncomplete
}

// selectContext 在预算内去重、扩展父块并分配稳定引用 ID。
func selectContext(ctx context.Context, counter tokenCounter, candidates []contextCandidate, budget int) ([]selectedContext, error) {
	// TODO 4：按相关度在预算内选择，并分配稳定 Citation ID。
	return nil, errExerciseIncomplete
}

// validateSelectedCitations 确保回答引用只指向已选证据。
func validateSelectedCitations(selected []selectedContext, citations []string) error {
	// TODO 5：校验文件名、页码和 Citation ID 都来自实际选中上下文。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“上下文预算与引用”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return validateSelectedCitations(nil, nil)
}
