package main

import (
	"context"
	"errors"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type goldenCase struct {
	ID               string   `json:"id"`
	Query            string   `json:"query"`
	RelevantChunkIDs []string `json:"relevant_chunk_ids"`
	ExpectedFacts    []string `json:"expected_facts"`
	ExpectedTools    []string `json:"expected_tools"`
	Tags             []string `json:"tags"`
}

// loadGoldenDataset 按行严格解析 JSONL，并拒绝重复 ID、未知字段和空目标。
func loadGoldenDataset(path string) ([]goldenCase, error) {
	return nil, errExerciseIncomplete
}

// validateGoldenCase 校验单条评估样本的必填字段和目标约束。
func validateGoldenCase(value goldenCase) error {
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“Golden Dataset”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：定义版本化 JSONL Schema 和稳定 Case ID。
	// TODO 2：实现 loadGoldenDataset，限制文件大小并逐行严格解码。
	// TODO 3：实现 validateGoldenCase，拒绝重复 ID、空 Query 和无目标样本。
	// TODO 4：按主题、难度、租户和风险标签组织 fixture。
	// TODO 5：记录数据集版本与变更原因，保证评估可重复。
	return errExerciseIncomplete
}
