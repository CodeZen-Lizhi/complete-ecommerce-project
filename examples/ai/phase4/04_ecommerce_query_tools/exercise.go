package main

import (
	"context"
	"errors"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

// TODO 1：保持商品、库存、订单和物流接口狭窄，只返回只读 DTO。
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

// validateBusinessQuery 校验资源 ID、分页和返回大小。
func validateBusinessQuery(resourceID string, limit int) error {
	// TODO 2：拒绝空 ID、非法分页和超大返回，并禁止 Tool 直接访问 GORM。
	return errExerciseIncomplete
}

// queryProductTool 使用注入的 Service 查询有限商品字段。
func queryProductTool(ctx context.Context, service productQueryService, productID string) (productView, error) {
	// TODO 3：注入确定性 Fake Service，完成只读商品查询和错误传播。
	return productView{}, errExerciseIncomplete
}

// queryOrderTool 在可信身份边界内执行订单查询。
func queryOrderTool(ctx context.Context, service orderQueryService, orderID string) (orderView, error) {
	// TODO 4：从 Context 获取可信用户，并校验订单归属。
	return orderView{}, errExerciseIncomplete
}

// verifyBusinessQueryFailures 覆盖查询工具失败边界。
func verifyBusinessQueryFailures(ctx context.Context) error {
	// TODO 5：覆盖 NotFound、跨用户、依赖错误和超大返回场景。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“电商业务查询工具”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return verifyBusinessQueryFailures(ctx)
}
