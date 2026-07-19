package main

import (
	"context"
	"errors"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
)

// Judge 模型配置集中放在顶部，练习时直接替换占位值。
const (
	baseURL   = "http://localhost:8084/v1"
	apiKey    = "replace-with-your-api-key"
	modelName = "gpt-5.4-mini"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type answerCase struct {
	ID               string
	ExpectedFacts    []string
	AllowedCitations map[string]struct{}
	ShouldRefuse     bool
}

type judgeModelConfig struct {
	BaseURL string
	APIKey  string
	Model   string
}

type answerOutput struct {
	Text      string
	Citations []string
}

type judgeResult struct {
	Score       float64  `json:"score" jsonschema:"minimum=0,maximum=1"`
	Reasons     []string `json:"reasons" jsonschema:"minItems=1,minLength=1"`
	EvidenceIDs []string `json:"evidence_ids"`
}

var _ model.BaseChatModel = (*openai.ChatModel)(nil)

// loadAnswerCases 从 Golden Dataset 读取事实、引用和拒答标记。
func loadAnswerCases(path string) ([]answerCase, error) {
	// TODO 1：读取期望事实、允许引用和 ShouldRefuse，并拒绝重复 Case ID。
	return nil, errExerciseIncomplete
}

// validateCitations 确保引用只指向实际提供的证据。
func validateCitations(output answerOutput, allowed map[string]struct{}) error {
	// TODO 2：拒绝不存在或未发送给模型的证据 ID。
	return errExerciseIncomplete
}

// evaluateFacts 使用确定性规则检查期望事实和无依据拒答。
func evaluateFacts(input answerCase, output answerOutput) error {
	// TODO 3：检查事实覆盖，并验证无证据问题明确拒答。
	return errExerciseIncomplete
}

// newJudgeModel 使用真实 Eino OpenAI ChatModel 和严格 JSON Schema 创建可调用的 Judge。
func newJudgeModel(ctx context.Context, config judgeModelConfig) (model.BaseChatModel, error) {
	// TODO 4：调用真实模型，并用 judgeResult JSON Schema 严格解析评分。
	return nil, errExerciseIncomplete
}

// compareJudgeWithLabels 比较 Judge 评分和人工标签。
func compareJudgeWithLabels(ctx context.Context, config judgeModelConfig, datasetPath string) error {
	// TODO 5：输出一致率、失败事实、错误引用和人工抽样复核清单。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“回答事实与引用检查”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return compareJudgeWithLabels(ctx, judgeModelConfig{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Model:   modelName,
	}, "")
}
