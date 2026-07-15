package main

import (
	"context"
	"errors"
	"time"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type checkpoint struct {
	TaskID        string
	Version       int
	Step          string
	PendingAction string
	ExpiresAt     time.Time
	Status        string
}

type checkpointStore interface {
	// SaveInterrupt 原子保存中断状态和待确认动作。
	SaveInterrupt(ctx context.Context, value checkpoint) error
	// LoadForResume 读取可恢复且未过期的 Checkpoint。
	LoadForResume(ctx context.Context, taskID string) (checkpoint, error)
}

type confirmation struct {
	TaskID     string
	Version    int
	ApproverID string
	Approved   bool
}

// resumeFromCheckpoint 校验确认人与状态版本后幂等恢复。
func resumeFromCheckpoint(ctx context.Context, store checkpointStore, input confirmation) error {
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“Interrupt、Checkpoint 与 Resume”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：定义高风险节点、PendingAction 摘要和 Checkpoint 版本。
	// TODO 2：中断前使用 checkpointStore 保存状态，禁止提前执行副作用。
	// TODO 3：返回确认标识，并校验确认人、版本和过期时间。
	// TODO 4：实现 resumeFromCheckpoint，拒绝、过期和重复确认均幂等。
	// TODO 5：覆盖进程重启、版本冲突和确认后恢复路径。
	return errExerciseIncomplete
}
