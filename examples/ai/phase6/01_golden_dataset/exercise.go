package main

import (
	"context"
	"errors"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

// TODO 1：补全版本化 JSONL Schema，并约定稳定 Case ID。
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
	// TODO 2：限制文件大小，逐行严格解码并拒绝未知字段。
	return nil, errExerciseIncomplete
}

// validateGoldenCase 校验单条评估样本的必填字段和目标约束。
func validateGoldenCase(value goldenCase) error {
	// TODO 3：拒绝重复 ID、空 Query 和没有任何评估目标的样本。
	return errExerciseIncomplete
}

// organizeGoldenFixtures 按主题、难度、租户和风险标签组织样本。
func organizeGoldenFixtures(cases []goldenCase) (map[string][]goldenCase, error) {
	// TODO 4：校验标签集合并稳定分组，不能丢失未标记样本。
	return nil, errExerciseIncomplete
}

// recordDatasetVersion 保存数据集版本和变更原因。
func recordDatasetVersion(ctx context.Context, version string, reason string) error {
	// TODO 5：保证同一版本可重复读取，并拒绝空变更原因。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“Golden Dataset”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return recordDatasetVersion(ctx, "", "")
}
