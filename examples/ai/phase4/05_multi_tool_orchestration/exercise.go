package main

import (
	"context"
	"errors"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type toolCall struct {
	ID        string
	Name      string
	Arguments string
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

// bindRealTools 将真实 Eino 工具同时绑定到 ToolCallingChatModel 和 ToolsNode。
func bindRealTools(ctx context.Context, chatModel model.ToolCallingChatModel, tools []tool.BaseTool) (model.ToolCallingChatModel, *compose.ToolsNode, error) {
	return nil, nil, errExerciseIncomplete
}

// runToolLoop 执行有限工具循环，并检测未知、重复和无调用场景。
func runToolLoop(ctx context.Context, model modelTurn, registry toolRegistry, maxCalls int) (string, error) {
	return "", errExerciseIncomplete
}

// runExercise 按执行顺序组织“多工具注册与连续调用”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：为工具提供唯一名称、明确描述和不重叠的职责。
	// TODO 2：实现 bindRealTools，通过模型 WithTools 和 compose.NewToolNode 不可变绑定真实工具，拒绝重复名称。
	// TODO 3：真实调用 ToolCallingChatModel，处理直接回答、单次调用和连续调用，再把 Tool Result 返回模型。
	// TODO 4：严格关联 Tool Call ID 与 toolOutput，未知工具立即失败。
	// TODO 5：限制总调用次数并检测重复调用指纹。
	return errExerciseIncomplete
}
