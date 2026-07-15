package main

import (
	"context"
	"errors"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type secretConfig struct {
	APIKey      string
	DatabaseDSN string
}

type logFields map[string]string

// loadSecretConfig 从环境读取必需配置，不提供真实默认值。
func loadSecretConfig() (secretConfig, error) {
	return secretConfig{}, errExerciseIncomplete
}

// redactFields 按字段白名单输出日志，并脱敏 Token、Cookie、密码和 PII。
func redactFields(fields logFields) (logFields, error) {
	return nil, errExerciseIncomplete
}

// containsSecret 检查输出文本是否包含任一原始敏感值。
func containsSecret(output string, secrets []string) (bool, error) {
	return false, errExerciseIncomplete
}

// runExercise 按执行顺序组织“Secret 管理与日志脱敏”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：从环境读取 API Key 和数据库配置，拒绝源码真实默认值。
	// TODO 2：实现 loadSecretConfig，区分必需项、可选项和占位符。
	// TODO 3：实现 redactFields，只允许低敏字段并遮蔽 Token、Cookie、密码和 PII。
	// TODO 4：实现 containsSecret，用测试扫描日志和错误文本。
	// TODO 5：验证配置错误不回显原始 Secret，并记录最小诊断信息。
	return errExerciseIncomplete
}
