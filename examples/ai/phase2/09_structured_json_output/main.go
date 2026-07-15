package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/schema"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type structuredAnswer struct {
	Summary    string   `json:"summary"`
	KeyPoints  []string `json:"key_points"`
	Confidence float64  `json:"confidence"`
}

type modelClient interface {
	// Generate 根据消息生成原始文本；实现不得自动清洗或修复 JSON。
	Generate(ctx context.Context, messages []*schema.Message) (string, error)
}

type incompleteModel struct{}

// Generate 保持骨架离线，并提示学习者替换为真实 Eino 模型适配器。
func (incompleteModel) Generate(context.Context, []*schema.Message) (string, error) {
	return "", errExerciseIncomplete
}

// main 运行结构化 JSON Prompt、模型调用、解析和字段校验练习。
func main() {
	answer, err := runExercise(context.Background(), incompleteModel{}, "解释 Go interface")
	if err != nil {
		fmt.Printf("结构化输出练习失败: %v\n", err)
		return
	}

	fmt.Printf("摘要: %s\n关键点: %v\n置信度: %.2f\n", answer.Summary, answer.KeyPoints, answer.Confidence)
}

// buildStructuredMessages 构造要求模型只返回固定 JSON 结构的消息。
func buildStructuredMessages(ctx context.Context, question string) ([]*schema.Message, error) {
	if ctx == nil {
		return nil, errors.New("Context 不能为空")
	}
	if strings.TrimSpace(question) == "" {
		return nil, errors.New("问题不能为空")
	}

	// TODO 1：编写 System Message，明确禁止 Markdown 代码块和额外文本。
	// TODO 2：在 User Message 中给出 summary、key_points、confidence 的 JSON Schema 示例和当前问题。
	return nil, errExerciseIncomplete
}

// runExercise 串联消息构造、模型调用和严格 JSON 解析。
func runExercise(
	ctx context.Context,
	client modelClient,
	question string,
) (structuredAnswer, error) {
	if ctx == nil {
		return structuredAnswer{}, errors.New("Context 不能为空")
	}
	if client == nil {
		return structuredAnswer{}, errors.New("Model Client 不能为空")
	}

	messages, err := buildStructuredMessages(ctx, question)
	if err != nil {
		return structuredAnswer{}, fmt.Errorf("构造结构化消息失败: %w", err)
	}

	// TODO 3：调用 client.Generate；模型错误使用 %w 包装，不能把失败当成空 JSON。
	raw, err := client.Generate(ctx, messages)
	if err != nil {
		return structuredAnswer{}, fmt.Errorf("模型生成结构化输出失败: %w", err)
	}

	// TODO 4：把原始文本交给 decodeAndValidate；不要用字符串截取静默删除代码块或额外文本。
	answer, err := decodeAndValidate([]byte(raw))
	if err != nil {
		return structuredAnswer{}, fmt.Errorf("解析结构化输出失败: %w", err)
	}
	return answer, nil
}

// decodeAndValidate 解析模型 JSON，并校验所有业务字段。
func decodeAndValidate(raw []byte) (structuredAnswer, error) {
	if len(raw) == 0 {
		return structuredAnswer{}, errors.New("模型输出不能为空")
	}
	if !json.Valid(raw) {
		return structuredAnswer{}, errors.New("模型输出不是合法 JSON")
	}

	// TODO 5：使用 json.Decoder，并调用 DisallowUnknownFields 拒绝未知字段。
	// TODO 6：解码后确认没有第二个 JSON 值或尾随非空内容。
	// TODO 7：校验 Summary 非空、KeyPoints 非空且元素无空白、Confidence 位于 [0,1]。
	// TODO 8：错误中区分 JSON 语法、结构和字段业务校验，并使用 %w 保留底层错误。
	return structuredAnswer{}, errExerciseIncomplete
}
