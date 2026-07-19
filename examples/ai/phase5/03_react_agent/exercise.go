package main

import (
	"context"
	"errors"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type agentToolCall struct {
	ID        string
	Name      string
	Arguments string
}

type agentToolResult struct {
	ToolCallID string
	Content    string
}

type reactModel interface {
	// Next 根据消息决定最终回答或下一组工具调用。
	Next(ctx context.Context, results []agentToolResult) ([]agentToolCall, string, error)
}

type reactToolsNode interface {
	// Execute 执行工具调用并返回与调用 ID 对应的结果。
	Execute(ctx context.Context, calls []agentToolCall) ([]agentToolResult, error)
}

// newEinoReActAgent 使用真实 ToolCallingChatModel 和 ToolsNodeConfig 创建 Eino ReAct Agent。
func newEinoReActAgent(ctx context.Context, chatModel model.ToolCallingChatModel, tools []tool.BaseTool, maxSteps int) (*react.Agent, error) {
	_ = compose.ToolsNodeConfig{Tools: tools}
	return nil, errExerciseIncomplete
}

// runReAct 执行有最大步骤限制的模型-工具循环。
func runReAct(ctx context.Context, model reactModel, tools reactToolsNode, maxSteps int) (string, error) {
	return "", errExerciseIncomplete
}

// runExercise 按执行顺序组织“ReAct Agent”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：准备只读工具及严格 Schema，并以不可变方式绑定到模型。
	// TODO 2：实现 newEinoReActAgent，通过 react.NewAgent 和 compose.ToolsNodeConfig 创建真实 Agent。
	// TODO 3：实现 runReAct，关联 Tool Call ID 和 Tool Result。
	// TODO 4：处理模型直接回答、单次调用、连续调用和工具失败。
	// TODO 5：设置最大步骤并检测重复调用，超限返回可诊断错误。
	return errExerciseIncomplete
}
