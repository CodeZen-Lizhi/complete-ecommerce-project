package main

import (
	"context"
	"errors"

	"github.com/cloudwego/eino/compose"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type routeKind string

const (
	routeKnowledge routeKind = "knowledge"
	routeProduct   routeKind = "product"
	routeOrder     routeKind = "order"
	routeClarify   routeKind = "clarify"
)

type routeInput struct {
	Question string
	UserID   string
}

type routeOutput struct {
	Route  routeKind
	Answer string
}

// classifyRoute 使用确定性规则优先选择路由，无法判断时返回澄清。
func classifyRoute(input routeInput) (routeKind, error) {
	return "", errExerciseIncomplete
}

// buildRoutingGraph 注册分支节点、边和结束节点并编译 Graph。
func buildRoutingGraph(ctx context.Context) (compose.Runnable[routeInput, routeOutput], error) {
	return nil, errExerciseIncomplete
}

// runExercise 按执行顺序组织“Graph 条件路由”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：定义知识、商品、订单和澄清路由类型。
	// TODO 2：实现 classifyRoute，固定业务规则优先于模型判断。
	// TODO 3：使用 compose.NewGraph 注册路由节点、三个业务节点和澄清节点。
	// TODO 4：通过真实分支边连接结束节点并 Compile 为 compose.Runnable。
	// TODO 5：覆盖未知分类、节点错误、取消和每条分支可达性。
	return errExerciseIncomplete
}
