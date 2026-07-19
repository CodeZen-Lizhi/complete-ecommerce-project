package main

import (
	"context"
	"errors"

	"github.com/cloudwego/eino/compose"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

// TODO 1：补全知识、商品、订单和澄清路由类型及输入输出字段。
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
	// TODO 2：固定业务规则优先于模型判断，无法分类时返回 routeClarify。
	return "", errExerciseIncomplete
}

// newRoutingGraph 注册路由节点、业务节点和澄清节点。
func newRoutingGraph() (*compose.Graph[routeInput, routeOutput], error) {
	// TODO 3：使用 compose.NewGraph 注册路由节点、三个业务节点和澄清节点。
	return nil, errExerciseIncomplete
}

// compileRoutingGraph 连接分支边和结束节点并编译 Graph。
func compileRoutingGraph(ctx context.Context, graph *compose.Graph[routeInput, routeOutput]) (compose.Runnable[routeInput, routeOutput], error) {
	// TODO 4：通过真实分支边连接结束节点，并 Compile 为 compose.Runnable。
	return nil, errExerciseIncomplete
}

// buildRoutingGraph 构建并编译完整路由 Graph。
func buildRoutingGraph(ctx context.Context) (compose.Runnable[routeInput, routeOutput], error) {
	graph, err := newRoutingGraph()
	if err != nil {
		return nil, err
	}
	return compileRoutingGraph(ctx, graph)
}

// verifyGraphRoutes 验证每条分支和失败路径都可达。
func verifyGraphRoutes(ctx context.Context) error {
	// TODO 5：覆盖未知分类、节点错误、Context 取消和每条分支可达性。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“Graph 条件路由”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return verifyGraphRoutes(ctx)
}
