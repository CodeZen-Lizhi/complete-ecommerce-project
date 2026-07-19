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
	// TODO 1：只接受可信认证边界提供的身份，并校验空用户与租户。
	return nil, errExerciseIncomplete
}

// identityFromContext 读取可信身份；缺失或无效身份明确失败。
func identityFromContext(ctx context.Context) (identity, error) {
	// TODO 2：读取可信身份，不接受模型参数里的 User ID 作为替代。
	return identity{}, errExerciseIncomplete
}

// authorizeResource 校验租户、资源归属和动作权限。
func authorizeResource(value identity, target resource, action string) error {
	// TODO 3：先校验租户，再校验资源归属和角色动作白名单。
	return errExerciseIncomplete
}

// classifyAuthorizationFailure 区分认证、授权和资源隐藏错误。
func classifyAuthorizationFailure(err error) (string, error) {
	// TODO 4：避免通过错误信息泄露跨租户资源是否存在。
	return "", errExerciseIncomplete
}

// recordAuthorizationAudit 记录最小授权审计字段。
func recordAuthorizationAudit(ctx context.Context, value identity, target resource, action string) error {
	// TODO 5：不得记录 Token、Cookie 或完整工具参数。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“Context 身份与再次授权”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return recordAuthorizationAudit(ctx, identity{}, resource{}, "")
}
