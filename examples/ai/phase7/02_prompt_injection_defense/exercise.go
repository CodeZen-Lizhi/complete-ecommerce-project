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

// loadInjectionCases 准备正常、直接注入和恶意文档样本。
func loadInjectionCases() ([]injectionCase, error) {
	// TODO 1：为每类攻击定义期望阻断结果和安全对照组。
	return nil, errExerciseIncomplete
}

// evaluateInjection 把检索内容视为不可信数据，并校验越权或泄密意图。
func evaluateInjection(input injectionCase, authorization authorizationContext) (securityDecision, error) {
	// TODO 2：检索内容不能覆盖 System 规则或请求 Secret。
	return securityDecision{}, errExerciseIncomplete
}

// enforceRetrieverScope 在检索阶段强制租户和文档权限过滤。
func enforceRetrieverScope(ctx context.Context, authorization authorizationContext) error {
	// TODO 3：拒绝空租户、越权文档和模型生成的权限条件。
	return errExerciseIncomplete
}

// authorizeToolCall 在固定权限边界内校验工具名称和资源范围。
func authorizeToolCall(authorization authorizationContext, toolName string, resourceID string) error {
	// TODO 4：模型选择工具后仍需按允许工具和资源集合再次授权。
	return errExerciseIncomplete
}

// recordInjectionAudit 记录攻击类型和阻断原因。
func recordInjectionAudit(ctx context.Context, decision securityDecision) error {
	// TODO 5：不得记录完整恶意内容、Prompt 或 Secret。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“Prompt Injection 防御”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return recordInjectionAudit(ctx, securityDecision{})
}
