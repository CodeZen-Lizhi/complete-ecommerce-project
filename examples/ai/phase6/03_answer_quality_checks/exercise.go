package main

import (
	"context"
	"errors"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type answerCase struct {
	ID               string
	ExpectedFacts    []string
	AllowedCitations map[string]struct{}
	ShouldRefuse     bool
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

// newJudgeModel 使用真实 Eino OpenAI ChatModel 和严格 JSON Schema 创建可调用的 Judge。
func newJudgeModel(ctx context.Context) (model.BaseChatModel, error) {
	return nil, errExerciseIncomplete
}

// validateCitations 确保引用只指向实际提供的证据。
func validateCitations(output answerOutput, allowed map[string]struct{}) error {
	return errExerciseIncomplete
}

// evaluateFacts 使用确定性规则检查期望事实和无依据拒答。
func evaluateFacts(input answerCase, output answerOutput) error {
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“回答事实与引用检查”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：从数据集读取期望事实、允许引用和拒答标记。
	// TODO 2：实现 validateCitations，拒绝不存在或未发送给模型的证据 ID。
	// TODO 3：实现 evaluateFacts，检查事实覆盖和无证据问题的明确拒答。
	// TODO 4：实现 newJudgeModel，真实调用模型并用 judgeResult JSON Schema 解析评分，不得静默修复输出。
	// TODO 5：把 Judge 评分与人工标签比较，输出一致率、失败事实、错误引用和人工抽样复核清单。
	return errExerciseIncomplete
}
