package main

import (
	"context"
	"errors"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type limitKey struct {
	TenantID string
	UserID   string
	Resource string
}

type rateLimiter interface {
	// Allow 在 Context 有效期内等待许可或返回拒绝。
	Allow(ctx context.Context, key limitKey) error
}

type concurrencyLimiter interface {
	// Acquire 获取有界并发许可；返回的 release 必须调用一次。
	Acquire(ctx context.Context) (release func(), err error)
}

// runLimited 在限速和并发上限内执行一次下游调用。
func runLimited(ctx context.Context, rate rateLimiter, concurrency concurrencyLimiter, key limitKey, call func(context.Context) error) error {
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“限流与并发控制”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：定义请求级、用户级、租户级和下游级 limitKey。
	// TODO 2：实现 rateLimiter，等待许可时响应 Context 取消。
	// TODO 3：实现 concurrencyLimiter，使用有界许可并保证所有路径释放。
	// TODO 4：实现 runLimited，区分排队超时、主动取消和限额拒绝。
	// TODO 5：压测突发流量并记录 P95、排队时间和拒绝率。
	return errExerciseIncomplete
}
