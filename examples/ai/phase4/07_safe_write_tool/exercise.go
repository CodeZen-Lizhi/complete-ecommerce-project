package main

import (
	"context"
	"errors"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type writeCommand struct {
	Action         string
	ResourceID     string
	IdempotencyKey string
	ConfirmationID string
	DryRun         bool
}

type writePreview struct {
	Summary string
	Amount  int64
}

type confirmationStore interface {
	// Verify 校验确认令牌、确认人、动作摘要和过期时间。
	Verify(ctx context.Context, command writeCommand, preview writePreview) error
}

type idempotencyStore interface {
	// Load 返回已有写操作结果和是否命中。
	Load(ctx context.Context, key string) (string, bool, error)
	// Save 原子记录写操作结果。
	Save(ctx context.Context, key string, result string) error
}

// executeWriteTool 在确认和幂等边界内执行一次受控写操作。
func executeWriteTool(ctx context.Context, command writeCommand, confirmations confirmationStore, idempotency idempotencyStore) (string, error) {
	return "", errExerciseIncomplete
}

// runExercise 按执行顺序组织“安全写操作模拟”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：默认以 DryRun 生成 writePreview，不立即执行副作用。
	// TODO 2：要求应用侧 ConfirmationID，模型文本不能替代用户确认。
	// TODO 3：使用 idempotencyStore 检查重复取消或退款。
	// TODO 4：再次校验可信身份、资源状态、金额和动作白名单。
	// TODO 5：执行后保存幂等结果与审计事件，验证重复调用不产生第二次副作用。
	return errExerciseIncomplete
}
