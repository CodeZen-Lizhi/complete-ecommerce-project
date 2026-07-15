package main

import (
	"context"
	"errors"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type inventoryQuery struct {
	ProductID string `json:"product_id"`
	Warehouse string `json:"warehouse"`
	Limit     int    `json:"limit"`
}

// decodeInventoryQuery 使用严格 JSON 解码并拒绝未知字段和尾随内容。
func decodeInventoryQuery(argumentsJSON string) (inventoryQuery, error) {
	return inventoryQuery{}, errExerciseIncomplete
}

// validateInventoryQuery 校验必填字段、枚举、范围和字符串长度。
func validateInventoryQuery(query inventoryQuery) error {
	return errExerciseIncomplete
}

type validationError struct {
	Field   string
	Message string
}

// Error 返回不包含原始恶意输入的字段级校验错误。
func (e validationError) Error() string {
	return e.Field + ": " + e.Message
}

// runExercise 按执行顺序组织“Tool 参数严格校验”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：实现 decodeInventoryQuery，使用 DisallowUnknownFields 并拒绝第二个 JSON 值。
	// TODO 2：实现 validateInventoryQuery，校验商品 ID、仓库白名单和 Limit 上限。
	// TODO 3：区分 JSON 语法、结构和业务校验错误。
	// TODO 4：把内部错误映射为有限的 validationError，不回显完整恶意输入。
	// TODO 5：使用表驱动测试覆盖空值、未知字段、越界和合法参数。
	return errExerciseIncomplete
}
