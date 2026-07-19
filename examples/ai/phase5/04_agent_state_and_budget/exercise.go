package main

import (
	"context"
	"errors"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

// TODO 1：补全可序列化状态机、乐观版本和预算字段。
type agentState struct {
	TaskID      string
	Version     int
	Step        int
	MaxSteps    int
	TokenUsed   int
	TokenBudget int
	CostUsed    float64
	CostBudget  float64
	Status      string
	LastError   string
}

type stateStore interface {
	// Load 读取任务的最新版本状态。
	Load(ctx context.Context, taskID string) (agentState, error)
	// Save 使用期望版本保存新状态，版本冲突明确失败。
	Save(ctx context.Context, expectedVersion int, state agentState) error
}

// checkBudget 在每一步前校验步骤、Token 和成本预算。
func checkBudget(state agentState) error {
	// TODO 2：在每一步前检查最大步骤、Token 和成本预算。
	return errExerciseIncomplete
}

// advanceState 生成下一版本状态，不修改传入值。
func advanceState(state agentState, tokens int, cost float64) (agentState, error) {
	// TODO 3：使用不可变更新记录步骤、工具结果和错误摘要。
	return agentState{}, errExerciseIncomplete
}

// persistFinalState 保存成功、失败或取消后的任务终态。
func persistFinalState(ctx context.Context, store stateStore, state agentState) error {
	// TODO 4：使用期望版本保存，并明确处理版本冲突。
	return errExerciseIncomplete
}

// verifyAgentStateBoundaries 覆盖状态和预算边界。
func verifyAgentStateBoundaries(ctx context.Context) error {
	// TODO 5：测试版本冲突、预算耗尽、空 TaskID 和状态恢复。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“Agent 状态与预算”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return verifyAgentStateBoundaries(ctx)
}
