package main

import (
	"context"
	"errors"
	"time"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type hnswConfig struct {
	M              int
	EFConstruction int
	EFSearch       int
	TopK           int
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

// evaluateHNSW 运行固定查询集并汇总召回、延迟和索引大小。
func evaluateHNSW(ctx context.Context, index hnswIndex, config hnswConfig, samples []searchSample) (hnswReport, error) {
	return hnswReport{}, errExerciseIncomplete
}

// runExercise 按执行顺序组织“HNSW 参数调优”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：校验 M、efConstruction、efSearch 和 TopK 均为正数。
	// TODO 2：准备固定语料、查询集和相关文档集合。
	// TODO 3：为每组参数构建索引，并通过 evaluateHNSW 执行相同查询。
	// TODO 4：计算 Recall@K、P50/P95 和索引大小，记录失败查询。
	// TODO 5：输出参数对比表，选择满足质量目标的最低延迟配置。
	return errExerciseIncomplete
}
