package main

import (
	"context"
	"errors"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type queryFilter struct {
	TenantID     string
	AllowedDocs  []string
	DocumentType string
	SinceUnix    int64
}

type queryRewriter interface {
	// Rewrite 把依赖历史的追问改写成独立查询。
	Rewrite(ctx context.Context, history []string, question string) (string, error)
	// MultiQuery 基于独立查询生成有限数量的等价查询。
	MultiQuery(ctx context.Context, query string, limit int) ([]string, error)
}

type filteredRetriever interface {
	// Retrieve 在强制 Metadata Filter 下执行检索。
	Retrieve(ctx context.Context, query string, filter queryFilter, topK int) ([]string, error)
}

// validateQueryFilter 校验租户、权限范围和查询数量边界。
func validateQueryFilter(filter queryFilter) error {
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“问题改写、Multi-Query 与 Metadata Filter”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：实现 validateQueryFilter，租户和权限范围不能为空且不能由模型覆盖。
	// TODO 2：使用 queryRewriter 把多轮追问改写为独立查询。
	// TODO 3：生成有上限的 Multi-Query，去重并限制并发。
	// TODO 4：每一路 filteredRetriever 调用都强制传入相同 Metadata Filter。
	// TODO 5：融合结果并记录每个查询的命中贡献、错误和耗时。
	return errExerciseIncomplete
}
