package main

import (
	"context"
	"errors"
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
	return embeddingInput{}, errExerciseIncomplete
}

// loadEmbeddingConfig 从环境变量读取 OpenAI-compatible Embedding 配置。
func loadEmbeddingConfig() (embeddingConfig, error) {
	return embeddingConfig{}, errExerciseIncomplete
}

// cosineSimilarity 计算两个等维非零向量的余弦相似度。
func cosineSimilarity(left []float64, right []float64) (float64, error) {
	return 0, errExerciseIncomplete
}

// rankTopK 按相似度稳定排序文本，并返回最多 topK 个结果。
func rankTopK(texts []string, vectors [][]float64, queryVector []float64, topK int) ([]scoredText, error) {
	return nil, errExerciseIncomplete
}

// runExercise 按执行顺序组织“Embedding、余弦相似度与 Top-K”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：实现 parseEmbeddingInput 和 loadEmbeddingConfig，先校验 CLI 与占位配置。
	// TODO 2：通过 Eino OpenAI-compatible Embedder 为候选文本和查询分别调用 EmbedStrings。
	// TODO 3：实现 cosineSimilarity，拒绝空向量、零向量和维度不一致。
	// TODO 4：实现 rankTopK，校验文本与向量数量一致并稳定处理并列分数。
	// TODO 5：打印向量维度、Top-K 文本、分数和耗时，不输出 API Key。
	return errExerciseIncomplete
}
