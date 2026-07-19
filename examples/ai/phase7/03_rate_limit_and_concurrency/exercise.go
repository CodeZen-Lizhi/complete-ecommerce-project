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

// buildLimitKey 定义请求、用户、租户和下游资源维度。
func buildLimitKey(tenantID string, userID string, resource string) (limitKey, error) {
	// TODO 1：拒绝空租户与资源，并明确匿名用户的稳定 Key。
	return limitKey{}, errExerciseIncomplete
}

// waitRateLimit 等待速率许可并响应 Context 取消。
func waitRateLimit(ctx context.Context, limiter rateLimiter, key limitKey) error {
	// TODO 2：区分正常许可、排队超时、主动取消和限额拒绝。
	return errExerciseIncomplete
}

// acquireConcurrencySlot 获取有界并发许可并返回释放函数。
func acquireConcurrencySlot(ctx context.Context, limiter concurrencyLimiter) (func(), error) {
	// TODO 3：保证成功获取后的所有路径只释放一次。
	return nil, errExerciseIncomplete
}

// runLimited 在限速和并发上限内执行一次下游调用。
func runLimited(ctx context.Context, rate rateLimiter, concurrency concurrencyLimiter, key limitKey, call func(context.Context) error) error {
	// TODO 4：按限速、并发、下游调用顺序组合，并保留原始错误原因。
	return errExerciseIncomplete
}

// benchmarkLimitedCalls 压测突发流量并记录关键指标。
func benchmarkLimitedCalls(ctx context.Context) error {
	// TODO 5：记录 P95、排队时间和拒绝率，检查并发上限没有泄漏。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“限流与并发控制”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return benchmarkLimitedCalls(ctx)
}
