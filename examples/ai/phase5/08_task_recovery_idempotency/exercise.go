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

// recoverTasks 有界恢复任务，并在副作用前查询幂等记录。
func recoverTasks(ctx context.Context, tasks recoveryStore, records idempotencyRecordStore, limit int) error {
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“任务恢复与幂等”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：启动时通过 recoveryStore 有界扫描可恢复任务。
	// TODO 2：从最近 Checkpoint 重建状态，损坏或缺失 Checkpoint 明确失败。
	// TODO 3：为每个副作用生成稳定幂等键并先查询 idempotencyRecordStore。
	// TODO 4：实现 recoverTasks，已成功副作用直接复用记录。
	// TODO 5：在多个崩溃点测试恢复，确认不会重复执行副作用。
	return errExerciseIncomplete
}
