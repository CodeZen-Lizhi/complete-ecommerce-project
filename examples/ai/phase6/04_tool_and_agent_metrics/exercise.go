package main

import (
	"context"
	"errors"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type expectedExecution struct {
	CaseID         string
	ExpectedTools  []string
	ExpectedParams map[string]map[string]string
	ExpectedStatus string
}

type actualExecution struct {
	CaseID         string
	ToolCalls      []string
	NormalizedArgs map[string]map[string]string
	Status         string
	Steps          int
	Termination    string
}

type executionMetrics struct {
	ToolSelectionRate  float64
	ParameterAccuracy  float64
	TaskCompletionRate float64
}

// evaluateExecution 对比期望工具、参数和 Agent 终态。
func evaluateExecution(expected []expectedExecution, actual []actualExecution) (executionMetrics, error) {
	return executionMetrics{}, errExerciseIncomplete
}

// runExercise 按执行顺序组织“Tool 与 Agent 指标”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：定义期望工具、规范化参数和任务终态。
	// TODO 2：校验 Case 对齐、重复 Tool Call 和未知参数。
	// TODO 3：实现 evaluateExecution，计算工具选择率和参数正确率。
	// TODO 4：统计 Agent 完成率、平均步骤、重复调用和终止原因。
	// TODO 5：输出失败轨迹并区分模型、工具和环境问题。
	return errExerciseIncomplete
}
