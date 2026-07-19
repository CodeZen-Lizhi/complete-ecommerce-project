package main

import (
	"context"
	"errors"
	"time"

	"github.com/cloudwego/eino/compose"
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

// compileInterruptibleGraph 使用真实 CheckPointStore 和中断节点配置编译 Eino Graph。
func compileInterruptibleGraph(ctx context.Context, store compose.CheckPointStore, interruptNode string) (compose.Runnable[string, string], error) {
	return nil, errExerciseIncomplete
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
	// TODO 2：实现 compileInterruptibleGraph，通过 compose.WithCheckPointStore 和 WithInterruptBeforeNodes 编译真实 Graph。
	// TODO 3：首次 Invoke 使用 compose.WithCheckPointID 保存中断状态，返回确认标识并校验确认人、版本和过期时间。
	// TODO 4：实现 resumeFromCheckpoint，拒绝、过期和重复确认均幂等。
	// TODO 5：使用相同 CheckPoint ID，并通过 compose.Resume 或 ResumeWithData 构造恢复 Context，观察进程重启、版本冲突和确认后路径。
	return errExerciseIncomplete
}
