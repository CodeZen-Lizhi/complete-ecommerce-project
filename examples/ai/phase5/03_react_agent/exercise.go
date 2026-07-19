package main

import (
	"context"
	"errors"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
)

// 模型配置集中放在顶部，练习时直接替换占位值。
const (
	baseURL   = "http://localhost:8084/v1"
	apiKey    = "replace-with-your-api-key"
	modelName = "gpt-5.4-mini"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type agentToolCall struct {
	ID        string
	Name      string
	Arguments string
}

type reactModelConfig struct {
	BaseURL string
	APIKey  string
	Model   string
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

// prepareReActTools 准备只读工具及严格参数 Schema。
func prepareReActTools(ctx context.Context) ([]tool.BaseTool, error) {
	// TODO 1：创建只读工具，并以不可变方式绑定严格 Schema。
	return nil, errExerciseIncomplete
}

// newEinoReActAgent 使用真实 ToolCallingChatModel 和 ToolsNodeConfig 创建 Eino ReAct Agent。
func newEinoReActAgent(ctx context.Context, chatModel model.ToolCallingChatModel, tools []tool.BaseTool, maxSteps int) (*react.Agent, error) {
	// TODO 2：通过 react.NewAgent 和 compose.ToolsNodeConfig 创建真实 Agent。
	_ = compose.ToolsNodeConfig{Tools: tools}
	return nil, errExerciseIncomplete
}

// runReAct 执行有最大步骤限制的模型-工具循环。
func runReAct(ctx context.Context, model reactModel, tools reactToolsNode, maxSteps int) (string, error) {
	// TODO 3：关联 Tool Call ID 与 Tool Result，并把结果交回模型继续推理。
	return "", errExerciseIncomplete
}

// handleReActTurn 处理直接回答、连续工具调用和工具失败。
func handleReActTurn(calls []agentToolCall, answer string, results []agentToolResult) (string, error) {
	// TODO 4：区分直接回答、单次调用、连续调用和工具失败。
	return "", errExerciseIncomplete
}

// runReActScenario 执行带步骤上限和重复检测的真实 Agent 场景。
func runReActScenario(ctx context.Context, config reactModelConfig) error {
	// TODO 5：设置最大步骤并检测重复调用，超限返回可诊断错误。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“ReAct Agent”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return runReActScenario(ctx, reactModelConfig{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Model:   modelName,
	})
}
