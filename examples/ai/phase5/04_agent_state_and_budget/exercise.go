package main

import (
	"context"
	"errors"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

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
	return errExerciseIncomplete
}

// advanceState 生成下一版本状态，不修改传入值。
func advanceState(state agentState, tokens int, cost float64) (agentState, error) {
	return agentState{}, errExerciseIncomplete
}

// runExercise 按执行顺序组织“Agent 状态与预算”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：定义可序列化 agentState 和乐观版本字段。
	// TODO 2：实现 checkBudget，在每步前检查最大步骤、Token 和成本。
	// TODO 3：实现 advanceState，使用不可变更新记录工具结果和错误摘要。
	// TODO 4：通过 stateStore 保存成功、失败和取消后的最终状态。
	// TODO 5：测试版本冲突、预算耗尽、空 TaskID 和状态恢复。
	return errExerciseIncomplete
}
