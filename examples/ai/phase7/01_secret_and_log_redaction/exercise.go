package main

import (
	"context"
	"errors"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

// TODO 1：声明必需与可选 Secret，但不得提供真实源码默认值。
type secretConfig struct {
	APIKey      string
	DatabaseDSN string
}

type logFields map[string]string

// loadSecretConfig 从环境读取必需配置，不提供真实默认值。
func loadSecretConfig() (secretConfig, error) {
	// TODO 2：从环境读取配置，并拒绝空值和已知占位符。
	return secretConfig{}, errExerciseIncomplete
}

// redactFields 按字段白名单输出日志，并脱敏 Token、Cookie、密码和 PII。
func redactFields(fields logFields) (logFields, error) {
	// TODO 3：只保留低敏白名单字段，并遮蔽 Token、Cookie、密码和 PII。
	return nil, errExerciseIncomplete
}

// containsSecret 检查输出文本是否包含任一原始敏感值。
func containsSecret(output string, secrets []string) (bool, error) {
	// TODO 4：扫描日志和错误文本，拒绝空 Secret 导致的误判。
	return false, errExerciseIncomplete
}

// verifySecretRedaction 验证配置错误和日志都不泄露原始值。
func verifySecretRedaction(ctx context.Context) error {
	// TODO 5：只记录最小诊断信息，并覆盖常见敏感字段变体。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“Secret 管理与日志脱敏”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return verifySecretRedaction(ctx)
}
