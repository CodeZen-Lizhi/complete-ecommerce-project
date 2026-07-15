package main

import (
	"context"
	"errors"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type identity struct {
	UserID   string
	TenantID string
	Roles    []string
}

type resource struct {
	ID       string
	TenantID string
	OwnerID  string
}

type identityContextKey struct{}

// withIdentity 把可信身份写入新的 Context。
func withIdentity(ctx context.Context, value identity) (context.Context, error) {
	return nil, errExerciseIncomplete
}

// identityFromContext 读取可信身份；缺失或无效身份明确失败。
func identityFromContext(ctx context.Context) (identity, error) {
	return identity{}, errExerciseIncomplete
}

// authorizeResource 校验租户、资源归属和动作权限。
func authorizeResource(value identity, target resource, action string) error {
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“Context 身份与再次授权”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：定义可信 identity 的来源，并实现 withIdentity 的空值校验。
	// TODO 2：实现 identityFromContext，不接受模型参数中的 user ID 作为可信身份。
	// TODO 3：实现 authorizeResource，先校验租户，再校验资源归属和角色动作。
	// TODO 4：区分未认证、无权限和资源不存在，避免泄露跨租户资源。
	// TODO 5：记录最小审计字段，不记录 Token、Cookie 或完整工具参数。
	return errExerciseIncomplete
}
