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
	// TODO 1：使用 DisallowUnknownFields，并拒绝第二个 JSON 值和尾随内容。
	return inventoryQuery{}, errExerciseIncomplete
}

// validateInventoryQuery 校验必填字段、枚举、范围和字符串长度。
func validateInventoryQuery(query inventoryQuery) error {
	// TODO 2：校验商品 ID、仓库白名单、Limit 上限和字符串长度。
	return errExerciseIncomplete
}

// classifyInventoryQueryError 区分语法、结构和业务字段错误。
func classifyInventoryQueryError(err error) (string, error) {
	// TODO 3：保留底层错误链，并返回稳定错误类别。
	return "", errExerciseIncomplete
}

type validationError struct {
	Field   string
	Message string
}

// Error 返回不包含原始恶意输入的字段级校验错误。
func (e validationError) Error() string {
	return e.Field + ": " + e.Message
}

// mapValidationError 把内部错误映射成有限字段级错误。
func mapValidationError(err error) (validationError, error) {
	// TODO 4：不得回显完整恶意输入或内部实现细节。
	return validationError{}, errExerciseIncomplete
}

// verifyInventoryQueryCases 覆盖参数解析和校验边界。
func verifyInventoryQueryCases(ctx context.Context) error {
	// TODO 5：使用表驱动样本覆盖空值、未知字段、越界和合法参数。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“Tool 参数严格校验”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return verifyInventoryQueryCases(ctx)
}
