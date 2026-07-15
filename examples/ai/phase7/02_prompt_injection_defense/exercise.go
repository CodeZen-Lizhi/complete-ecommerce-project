package main

import (
	"context"
	"errors"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type injectionCase struct {
	Name             string
	UserPrompt       string
	RetrievedContent string
	ExpectedBlocked  bool
}

type authorizationContext struct {
	TenantID      string
	AllowedDocIDs map[string]struct{}
	AllowedTools  map[string]struct{}
}

type securityDecision struct {
	Blocked bool
	Reason  string
}

// evaluateInjection 把检索内容视为不可信数据，并校验越权或泄密意图。
func evaluateInjection(input injectionCase, authorization authorizationContext) (securityDecision, error) {
	return securityDecision{}, errExerciseIncomplete
}

// authorizeToolCall 在固定权限边界内校验工具名称和资源范围。
func authorizeToolCall(authorization authorizationContext, toolName string, resourceID string) error {
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“Prompt Injection 防御”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：准备正常、直接注入和恶意文档三类 fixture。
	// TODO 2：实现 evaluateInjection，检索内容不能覆盖 System 规则或请求 Secret。
	// TODO 3：在 Retriever 阶段强制租户和文档权限过滤。
	// TODO 4：实现 authorizeToolCall，模型选择工具后仍需再次授权。
	// TODO 5：记录攻击类型和阻断原因，不记录完整恶意内容或 Secret。
	return errExerciseIncomplete
}
