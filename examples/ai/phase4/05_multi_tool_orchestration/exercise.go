package main

import (
	"context"
	"errors"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

// 模型配置集中放在顶部，练习时直接替换占位值。
const (
	baseURL   = "http://localhost:8084/v1"
	apiKey    = "replace-with-your-api-key"
	modelName = "gpt-5.4-mini"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type toolCall struct {
	ID        string
	Name      string
	Arguments string
}

type toolModelConfig struct {
	BaseURL string
	APIKey  string
	Model   string
}

type toolOutput struct {
	ToolCallID string
	Content    string
}

type toolRegistry interface {
	// Execute 根据名称执行工具，并返回与 Tool Call ID 对应的结果。
	Execute(ctx context.Context, call toolCall) (toolOutput, error)
}

type modelTurn interface {
	// Next 根据消息决定回答或生成下一组 Tool Call。
	Next(ctx context.Context, outputs []toolOutput) ([]toolCall, string, error)
}

// validateToolDefinitions 校验工具名称和职责没有冲突。
func validateToolDefinitions(ctx context.Context, tools []tool.BaseTool) error {
	// TODO 1：确保工具名称唯一、描述明确且职责不重叠。
	return errExerciseIncomplete
}

// bindRealTools 将真实 Eino 工具同时绑定到 ToolCallingChatModel 和 ToolsNode。
func bindRealTools(ctx context.Context, chatModel model.ToolCallingChatModel, tools []tool.BaseTool) (model.ToolCallingChatModel, *compose.ToolsNode, error) {
	// TODO 2：通过 WithTools 和 compose.NewToolNode 不可变绑定工具。
	return nil, nil, errExerciseIncomplete
}

// runToolLoop 执行有限工具循环，并检测未知、重复和无调用场景。
func runToolLoop(ctx context.Context, model modelTurn, registry toolRegistry, maxCalls int) (string, error) {
	// TODO 3：处理直接回答、单次调用和连续调用，并把 Tool Result 返回模型。
	return "", errExerciseIncomplete
}

// correlateToolOutput 校验工具结果与原 Tool Call 一一对应。
func correlateToolOutput(call toolCall, output toolOutput) error {
	// TODO 4：严格关联 Tool Call ID，遇到未知工具或错配结果立即失败。
	return errExerciseIncomplete
}

// runToolOrchestrationScenario 执行带预算和重复检测的完整场景。
func runToolOrchestrationScenario(ctx context.Context, config toolModelConfig) error {
	// TODO 5：限制总调用次数，并使用规范化参数生成重复调用指纹。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“多工具注册与连续调用”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return runToolOrchestrationScenario(ctx, toolModelConfig{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Model:   modelName,
	})
}
