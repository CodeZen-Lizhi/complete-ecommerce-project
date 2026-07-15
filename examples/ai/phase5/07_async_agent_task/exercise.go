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

// runWorker 执行单个任务并在成功、失败或取消时持久化终态。
func runWorker(ctx context.Context, repository taskRepository, workerID string) error {
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“异步 Agent 任务”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：定义任务状态机和合法状态转换。
	// TODO 2：通过 taskRepository.CreateOrGet 实现创建幂等。
	// TODO 3：实现 Claim 的原子领取和租约，避免多个 Worker 重复执行。
	// TODO 4：实现 runWorker，保存进度、Checkpoint、错误和取消状态。
	// TODO 5：模拟 Worker 崩溃和租约过期，验证任务可被重新领取。
	return errExerciseIncomplete
}
