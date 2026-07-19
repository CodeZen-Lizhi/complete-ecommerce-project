package main

import (
	"context"
	"errors"
	"time"
)

// Qdrant 配置集中放在顶部，练习时直接替换占位值。
const (
	qdrantBaseURL    = "http://localhost:6333"
	qdrantAPIKey     = "replace-with-your-qdrant-api-key"
	qdrantCollection = "ai_phase3_hnsw"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type hnswConfig struct {
	M              int
	EFConstruction int
	EFSearch       int
	TopK           int
}

type hnswBackendConfig struct {
	BaseURL    string
	APIKey     string
	Collection string
}

type searchSample struct {
	Query       string
	RelevantIDs map[string]struct{}
}

type hnswIndex interface {
	// Build 使用给定构建参数创建索引。
	Build(ctx context.Context, config hnswConfig) error
	// Search 使用给定搜索参数返回有序文档 ID。
	Search(ctx context.Context, query string, config hnswConfig) ([]string, error)
}

type hnswReport struct {
	RecallAtK  float64
	P95        time.Duration
	IndexBytes int64
}

// validateHNSWConfig 校验构建和搜索参数。
func validateHNSWConfig(config hnswConfig) error {
	// TODO 1：要求 M、EFConstruction、EFSearch 和 TopK 均为正数。
	return errExerciseIncomplete
}

// loadSearchSamples 准备固定语料和相关文档集合。
func loadSearchSamples() ([]searchSample, error) {
	// TODO 2：定义固定查询集和每个查询的 Relevant ID。
	return nil, errExerciseIncomplete
}

// buildHNSWVariant 使用一组参数构建真实索引。
func buildHNSWVariant(ctx context.Context, index hnswIndex, config hnswConfig) error {
	// TODO 3：构建索引，并确保不同参数组使用相同语料和查询。
	return errExerciseIncomplete
}

// evaluateHNSW 运行固定查询集并汇总召回、延迟和索引大小。
func evaluateHNSW(ctx context.Context, index hnswIndex, config hnswConfig, samples []searchSample) (hnswReport, error) {
	// TODO 4：计算 Recall@K、P50/P95、索引大小并记录失败查询。
	return hnswReport{}, errExerciseIncomplete
}

// compareHNSWReports 输出参数对比并选择满足质量目标的配置。
func compareHNSWReports(ctx context.Context, backend hnswBackendConfig) error {
	// TODO 5：输出参数对比表并选择满足质量目标的最低延迟配置。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“HNSW 参数调优”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return compareHNSWReports(ctx, hnswBackendConfig{
		BaseURL:    qdrantBaseURL,
		APIKey:     qdrantAPIKey,
		Collection: qdrantCollection,
	})
}
