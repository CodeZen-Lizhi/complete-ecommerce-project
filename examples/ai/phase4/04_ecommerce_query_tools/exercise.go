package main

import (
	"context"
	"errors"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type productView struct {
	ID    string
	Name  string
	Stock int
}

type orderView struct {
	ID        string
	UserID    string
	Status    string
	Logistics string
}

type productQueryService interface {
	// GetProduct 返回有限字段的商品与库存信息。
	GetProduct(ctx context.Context, productID string) (productView, error)
}

type orderQueryService interface {
	// GetOrder 返回属于指定用户的订单与物流信息。
	GetOrder(ctx context.Context, userID string, orderID string) (orderView, error)
}

// queryOrderTool 在可信身份边界内执行订单查询。
func queryOrderTool(ctx context.Context, service orderQueryService, orderID string) (orderView, error) {
	return orderView{}, errExerciseIncomplete
}

// runExercise 按执行顺序组织“电商业务查询工具”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：为商品、库存、订单和物流定义窄查询接口和只读 DTO。
	// TODO 2：校验 ID、分页和返回大小，禁止 Tool 直接访问 GORM。
	// TODO 3：注入 Fake Service，实现商品只读查询工具。
	// TODO 4：实现 queryOrderTool，从 Context 获取可信用户并校验订单归属。
	// TODO 5：覆盖 NotFound、跨用户、依赖错误和超大返回场景。
	return errExerciseIncomplete
}
