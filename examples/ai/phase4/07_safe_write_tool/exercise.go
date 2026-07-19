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

// buildWritePreview 默认以 DryRun 生成副作用预览。
func buildWritePreview(command writeCommand) (writePreview, error) {
	// TODO 1：返回动作、资源和金额摘要，不立即执行写操作。
	return writePreview{}, errExerciseIncomplete
}

// verifyWriteConfirmation 校验应用侧确认凭证。
func verifyWriteConfirmation(ctx context.Context, confirmations confirmationStore, command writeCommand, preview writePreview) error {
	// TODO 2：模型文本不能替代 ConfirmationID 和确认人校验。
	return errExerciseIncomplete
}

// loadIdempotentWriteResult 查询重复取消或退款的既有结果。
func loadIdempotentWriteResult(ctx context.Context, store idempotencyStore, key string) (string, bool, error) {
	// TODO 3：要求稳定幂等键，并在副作用前查询已有结果。
	return "", false, errExerciseIncomplete
}

// authorizeWriteCommand 再次校验身份、资源状态和动作白名单。
func authorizeWriteCommand(ctx context.Context, command writeCommand, preview writePreview) error {
	// TODO 4：拒绝越权资源、非法金额和非白名单动作。
	return errExerciseIncomplete
}

// executeWriteTool 在确认和幂等边界内执行一次受控写操作。
func executeWriteTool(ctx context.Context, command writeCommand, confirmations confirmationStore, idempotency idempotencyStore) (string, error) {
	// TODO 5：执行后原子保存幂等结果和审计事件，重复调用不得产生第二次副作用。
	return "", errExerciseIncomplete
}

// runExercise 按执行顺序组织“安全写操作模拟”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	_, err := executeWriteTool(ctx, writeCommand{}, nil, nil)
	return err
}
