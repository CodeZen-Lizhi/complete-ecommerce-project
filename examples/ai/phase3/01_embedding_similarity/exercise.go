package main

import (
	"context"
	"errors"
)

// 模型配置集中放在顶部，练习时直接替换占位值。
const (
	baseURL   = "http://localhost:8084/v1"
	apiKey    = "replace-with-your-api-key"
	modelName = "text-embedding-3-small"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type embeddingConfig struct {
	BaseURL string
	APIKey  string
	Model   string
}

type embeddingInput struct {
	Text       string
	Candidates []string
	TopK       int
}

type embedder interface {
	// EmbedStrings 为一组文本生成顺序一致的向量。
	EmbedStrings(ctx context.Context, texts []string) ([][]float64, error)
}

type scoredText struct {
	Text  string
	Score float64
}

// parseEmbeddingInput 解析 --text、--candidate 和 --top-k 命令行参数。
func parseEmbeddingInput(args []string) (embeddingInput, error) {
	// TODO 1：解析并校验 --text、--candidate 和 --top-k 参数。
	return embeddingInput{}, errExerciseIncomplete
}

// loadEmbeddingConfig 从顶部常量读取 OpenAI-compatible Embedding 配置。
func loadEmbeddingConfig() (embeddingConfig, error) {
	// TODO 2：读取顶部常量，并拒绝空值和占位 API Key。
	return embeddingConfig{BaseURL: baseURL, APIKey: apiKey, Model: modelName}, errExerciseIncomplete
}

// embedInputs 为查询和候选文本生成顺序一致的向量。
func embedInputs(ctx context.Context, embedder embedder, input embeddingInput) ([][]float64, []float64, error) {
	// TODO 3：调用 EmbedStrings，并校验响应数量和向量维度一致。
	return nil, nil, errExerciseIncomplete
}

// cosineSimilarity 计算两个等维非零向量的余弦相似度。
func cosineSimilarity(left []float64, right []float64) (float64, error) {
	// TODO 4：拒绝空向量、零向量和维度不一致，再计算余弦相似度。
	return 0, errExerciseIncomplete
}

// rankTopK 按相似度稳定排序文本，并返回最多 topK 个结果。
func rankTopK(texts []string, vectors [][]float64, queryVector []float64, topK int) ([]scoredText, error) {
	// TODO 5：校验文本与向量数量一致，并稳定处理并列分数。
	return nil, errExerciseIncomplete
}

// reportEmbeddingResults 输出维度、排名、分数和耗时。
func reportEmbeddingResults(ctx context.Context) error {
	// TODO 6：打印 Top-K 结果和耗时，但不得输出 API Key。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“Embedding、余弦相似度与 Top-K”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return reportEmbeddingResults(ctx)
}
