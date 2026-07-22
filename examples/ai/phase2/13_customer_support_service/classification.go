package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
)

// responseStyle 是回答阶段可使用的固定表达风格。
type responseStyle string

const (
	styleConcise    responseStyle = "concise"
	styleGuided     responseStyle = "guided"
	styleEmpathetic responseStyle = "empathetic"
)

// classification 是通过严格结构化输出校验后的分类结果。
type classification struct {
	Intent          customerIntent `json:"intent"`
	ResponseStyle   responseStyle  `json:"response_style"`
	RequiresHandoff bool           `json:"requires_handoff"`
}

// buildClassificationResponseFormat 构造分类阶段使用的严格原生 JSON Schema。
func buildClassificationResponseFormat() (*openai.ChatCompletionResponseFormat, error) {
	// TODO 10：基于分类结果构建 additionalProperties:false、必填字段和 Strict:true 的 JSON Schema。
	return nil, errExerciseIncomplete
}

// classificationSystemMessage 生成分类阶段的 System 规则，并明确禁止输出解释文本。
func classificationSystemMessage() (*schema.Message, error) {
	// TODO 3：使用 System Message 约束意图、风格、转人工字段及仅返回 Schema JSON 的职责。
	return nil, errExerciseIncomplete
}

// classificationFewShotMessages 返回三个固定的 User/Assistant Few-shot 分类样例。
func classificationFewShotMessages() ([]*schema.Message, error) {
	// TODO 4：添加商品建议、配送退换、售后升级三个角色正确的 Few-shot 消息对。
	return nil, errExerciseIncomplete
}

// buildClassificationMessages 按 System、Few-shot、当前 User 的顺序组装分类消息。
func buildClassificationMessages(message string) ([]*schema.Message, error) {
	if strings.TrimSpace(message) == "" {
		return nil, fmt.Errorf("分类消息不能为空")
	}
	systemMessage, err := classificationSystemMessage()
	if err != nil {
		return nil, err
	}
	fewShotMessages, err := classificationFewShotMessages()
	if err != nil {
		return nil, err
	}
	messages := make([]*schema.Message, 0, len(fewShotMessages)+2)
	messages = append(messages, systemMessage)
	messages = append(messages, fewShotMessages...)
	messages = append(messages, schema.UserMessage(message))
	return messages, nil
}

// classifyCustomerMessage 通过治理层执行分类调用，并严格解析模型原始 JSON。
func classifyCustomerMessage(
	ctx context.Context,
	provider chatProvider,
	governance governanceConfig,
	message string,
) (classification, stageMetrics, error) {
	if ctx == nil {
		return classification{}, stageMetrics{}, fmt.Errorf("分类 Context 不能为空")
	}
	if provider == nil {
		return classification{}, stageMetrics{}, fmt.Errorf("分类 Provider 不能为空")
	}
	messages, err := buildClassificationMessages(message)
	if err != nil {
		return classification{}, stageMetrics{}, err
	}
	result, metrics, err := governedGenerate(ctx, provider, governance, modelCall{Messages: messages})
	if err != nil {
		return classification{}, metrics, fmt.Errorf("模型分类失败: %w", err)
	}
	classificationResult, err := decodeAndValidateClassification([]byte(result.Content))
	if err != nil {
		return classification{}, metrics, fmt.Errorf("分类结果无效: %w", err)
	}
	return classificationResult, metrics, nil
}

// decodeAndValidateClassification 严格解码模型 JSON，并验证所有业务枚举与字段。
func decodeAndValidateClassification(raw []byte) (classification, error) {
	if len(raw) == 0 {
		return classification{}, fmt.Errorf("分类 JSON 不能为空")
	}

	// TODO 11：使用 DisallowUnknownFields 和第二次 Decode 拒绝尾随值，并校验 intent/style/布尔字段。
	return classification{}, errExerciseIncomplete
}
