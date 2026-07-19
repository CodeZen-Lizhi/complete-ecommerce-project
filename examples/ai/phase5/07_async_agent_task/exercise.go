package main

import (
	"context"
	"errors"
	"time"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type taskStatus string

const (
	taskPending   taskStatus = "pending"
	taskRunning   taskStatus = "running"
	taskSucceeded taskStatus = "succeeded"
	taskFailed    taskStatus = "failed"
	taskCanceled  taskStatus = "canceled"
)

type agentTask struct {
	ID             string
	IdempotencyKey string
	Status         taskStatus
	Progress       int
	LeaseUntil     time.Time
	LastError      string
}

type taskRepository interface {
	// CreateOrGet 按幂等键创建任务或返回已有任务。
	CreateOrGet(ctx context.Context, task agentTask) (agentTask, error)
	// Claim 原子领取可运行任务并设置租约。
	Claim(ctx context.Context, workerID string, lease time.Duration) (agentTask, error)
	// SaveProgress 保存合法状态转换和进度。
	SaveProgress(ctx context.Context, task agentTask) error
}

// validateTaskTransition 校验任务状态机的合法转换。
func validateTaskTransition(from taskStatus, to taskStatus) error {
	// TODO 1：定义 Pending、Running、Succeeded、Failed、Canceled 的合法边。
	return errExerciseIncomplete
}

// createIdempotentTask 按幂等键创建或读取任务。
func createIdempotentTask(ctx context.Context, repository taskRepository, task agentTask) (agentTask, error) {
	// TODO 2：调用 CreateOrGet，并校验命中任务与输入语义一致。
	return agentTask{}, errExerciseIncomplete
}

// claimAgentTask 原子领取任务并设置租约。
func claimAgentTask(ctx context.Context, repository taskRepository, workerID string, lease time.Duration) (agentTask, error) {
	// TODO 3：避免多个 Worker 同时领取，并拒绝非法租约。
	return agentTask{}, errExerciseIncomplete
}

// runWorker 执行单个任务并在成功、失败或取消时持久化终态。
func runWorker(ctx context.Context, repository taskRepository, workerID string) error {
	// TODO 4：保存进度、Checkpoint、错误和取消状态。
	return errExerciseIncomplete
}

// verifyTaskRecovery 模拟 Worker 崩溃和租约过期。
func verifyTaskRecovery(ctx context.Context) error {
	// TODO 5：确认任务可以重新领取且不会重复执行已完成步骤。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“异步 Agent 任务”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return verifyTaskRecovery(ctx)
}
