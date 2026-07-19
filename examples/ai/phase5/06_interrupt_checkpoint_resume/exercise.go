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

// newPendingCheckpoint 创建高风险动作的待确认状态。
func newPendingCheckpoint(taskID string, action string, expiresAt time.Time) (checkpoint, error) {
	// TODO 1：填写 PendingAction 摘要、初始版本、过期时间和状态。
	return checkpoint{}, errExerciseIncomplete
}

// compileInterruptibleGraph 使用真实 CheckPointStore 和中断节点配置编译 Eino Graph。
func compileInterruptibleGraph(ctx context.Context, store compose.CheckPointStore, interruptNode string) (compose.Runnable[string, string], error) {
	// TODO 2：通过 WithCheckPointStore 和 WithInterruptBeforeNodes 编译真实 Graph。
	return nil, errExerciseIncomplete
}

// invokeUntilInterrupt 首次调用 Graph 并保存中断状态。
func invokeUntilInterrupt(ctx context.Context, graph compose.Runnable[string, string], checkpointID string) error {
	// TODO 3：通过 WithCheckPointID 调用 Graph，返回确认标识并校验版本和过期时间。
	return errExerciseIncomplete
}

// resumeFromCheckpoint 校验确认人与状态版本后幂等恢复。
func resumeFromCheckpoint(ctx context.Context, store checkpointStore, input confirmation) error {
	// TODO 4：让拒绝、过期、版本冲突和重复确认都保持幂等。
	return errExerciseIncomplete
}

// resumeEinoGraph 使用同一 Checkpoint ID 恢复真实 Eino Graph。
func resumeEinoGraph(ctx context.Context) error {
	// TODO 5：通过 compose.Resume 或 ResumeWithData 恢复，并观察进程重启后的路径。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“Interrupt、Checkpoint 与 Resume”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return resumeEinoGraph(ctx)
}
