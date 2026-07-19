package main

import (
	"context"
	"errors"
)

// 模型与向量库配置集中放在顶部，练习时直接替换占位值。
const (
	modelBaseURL     = "http://localhost:8084/v1"
	modelAPIKey      = "replace-with-your-api-key"
	modelName        = "gpt-5.4-mini"
	qdrantBaseURL    = "http://localhost:6333"
	qdrantAPIKey     = "replace-with-your-qdrant-api-key"
	qdrantCollection = "ai_phase3_filtered"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type queryFilter struct {
	TenantID     string
	AllowedDocs  []string
	DocumentType string
	SinceUnix    int64
}

type queryBackendConfig struct {
	ModelBaseURL     string
	ModelAPIKey      string
	ModelName        string
	QdrantBaseURL    string
	QdrantAPIKey     string
	QdrantCollection string
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
	// TODO 1：租户和权限范围不能为空，并且不能被模型输出覆盖。
	return errExerciseIncomplete
}

// rewriteStandaloneQuery 把多轮追问改写成独立查询。
func rewriteStandaloneQuery(ctx context.Context, rewriter queryRewriter, history []string, question string) (string, error) {
	// TODO 2：调用 Rewrite，并拒绝空问题或空改写结果。
	return "", errExerciseIncomplete
}

// buildMultiQueries 生成有上限且去重的等价查询。
func buildMultiQueries(ctx context.Context, rewriter queryRewriter, query string, limit int) ([]string, error) {
	// TODO 3：限制查询数量、去重并为后续并发设置上限。
	return nil, errExerciseIncomplete
}

// retrieveWithFilter 为每一路查询强制应用相同 Metadata Filter。
func retrieveWithFilter(ctx context.Context, retriever filteredRetriever, queries []string, filter queryFilter, topK int) ([][]string, error) {
	// TODO 4：每次 Retrieve 都传入同一租户与权限过滤条件。
	return nil, errExerciseIncomplete
}

// fuseQueryResults 汇总多路命中贡献、错误和耗时。
func fuseQueryResults(ctx context.Context, backend queryBackendConfig) error {
	// TODO 5：融合结果并输出每个查询的命中贡献、错误和耗时。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“问题改写、Multi-Query 与 Metadata Filter”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return fuseQueryResults(ctx, queryBackendConfig{
		ModelBaseURL:     modelBaseURL,
		ModelAPIKey:      modelAPIKey,
		ModelName:        modelName,
		QdrantBaseURL:    qdrantBaseURL,
		QdrantAPIKey:     qdrantAPIKey,
		QdrantCollection: qdrantCollection,
	})
}
