package main

import (
	"context"
	"errors"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type toolInfo struct {
	Name        string
	Description string
	InputSchema []byte
}

type readonlyTool interface {
	// Info 返回供模型选择工具的名称、描述和参数 Schema。
	Info(ctx context.Context) (toolInfo, error)
	// Invoke 校验 JSON 参数并执行无副作用查询。
	Invoke(ctx context.Context, argumentsJSON string) (string, error)
}

type toolExecutor interface {
	// Execute 根据工具名调用已注册的只读工具。
	Execute(ctx context.Context, toolName string, argumentsJSON string) (string, error)
}

// newEinoToolsNode 使用真实 Eino BaseTool 列表创建 ToolsNode；自定义执行器不能替代此步骤。
func newEinoToolsNode(ctx context.Context, tools []tool.BaseTool) (*compose.ToolsNode, error) {
	return nil, errExerciseIncomplete
}

// validateToolInfo 校验工具名、描述和 JSON Schema。
func validateToolInfo(info toolInfo) error {
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“本地只读工具与 ToolsNode”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：定义一个用途单一的本地只读工具及严格输入 Schema。
	// TODO 2：实现 validateToolInfo，拒绝空名称、模糊描述和非法 Schema。
	// TODO 3：实现 readonlyTool.Invoke，严格解析参数并传播 Context 取消。
	// TODO 4：实现 newEinoToolsNode，通过 compose.NewToolNode 注册真实工具并执行一次 Tool Call。
	// TODO 5：限制返回大小，错误不泄露内部数据或敏感字段。
	return errExerciseIncomplete
}
