package main

import (
	"context"
	"errors"
	"time"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type recoveryTask struct {
	ID            string
	CheckpointRef string
	Status        string
	LeaseUntil    time.Time
}

type recoveryStore interface {
	// ListRecoverable 返回租约过期或处于可恢复状态的任务。
	ListRecoverable(ctx context.Context, limit int) ([]recoveryTask, error)
	// LoadCheckpoint 读取任务最近的可用 Checkpoint。
	LoadCheckpoint(ctx context.Context, reference string) ([]byte, error)
}

type idempotencyRecordStore interface {
	// LoadResult 查询副作用是否已经成功执行。
	LoadResult(ctx context.Context, key string) (string, bool, error)
	// SaveResult 原子记录副作用结果。
	SaveResult(ctx context.Context, key string, result string) error
}

// listRecoveryBatch 有界扫描可恢复任务。
func listRecoveryBatch(ctx context.Context, store recoveryStore, limit int) ([]recoveryTask, error) {
	// TODO 1：调用 ListRecoverable，并拒绝非正或过大的扫描上限。
	return nil, errExerciseIncomplete
}

// restoreRecoveryTask 从最近 Checkpoint 重建任务状态。
func restoreRecoveryTask(ctx context.Context, store recoveryStore, task recoveryTask) error {
	// TODO 2：损坏或缺失 Checkpoint 必须明确失败。
	return errExerciseIncomplete
}

// buildSideEffectKey 为副作用生成稳定幂等键。
func buildSideEffectKey(task recoveryTask, operation string) (string, error) {
	// TODO 3：组合任务、步骤和操作标识，并拒绝空字段。
	return "", errExerciseIncomplete
}

// recoverTasks 有界恢复任务，并在副作用前查询幂等记录。
func recoverTasks(ctx context.Context, tasks recoveryStore, records idempotencyRecordStore, limit int) error {
	// TODO 4：已成功副作用直接复用记录，未命中时执行并原子保存结果。
	return errExerciseIncomplete
}

// verifyRecoveryCrashPoints 在多个崩溃点验证恢复幂等性。
func verifyRecoveryCrashPoints(ctx context.Context) error {
	// TODO 5：确认重启恢复不会重复执行任何已完成副作用。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“任务恢复与幂等”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return verifyRecoveryCrashPoints(ctx)
}
