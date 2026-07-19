package main

import (
	"context"
	"errors"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

// TODO 1：补全期望工具、规范化参数和任务终态契约。
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

// validateExecutionPairs 校验 Case、Tool Call 和参数对齐。
func validateExecutionPairs(expected []expectedExecution, actual []actualExecution) error {
	// TODO 2：拒绝缺失 Case、重复 Tool Call 和未知参数。
	return errExerciseIncomplete
}

// evaluateExecution 对比期望工具、参数和 Agent 终态。
func evaluateExecution(expected []expectedExecution, actual []actualExecution) (executionMetrics, error) {
	// TODO 3：计算工具选择率、参数正确率和任务完成率。
	return executionMetrics{}, errExerciseIncomplete
}

// summarizeAgentExecution 汇总步骤、重复调用和终止原因。
func summarizeAgentExecution(actual []actualExecution) error {
	// TODO 4：统计平均步骤、重复调用次数和各类终止原因。
	return errExerciseIncomplete
}

// reportExecutionFailures 输出失败轨迹及责任分类。
func reportExecutionFailures(ctx context.Context) error {
	// TODO 5：区分模型、工具和环境问题，并保留可复现轨迹。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“Tool 与 Agent 指标”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return reportExecutionFailures(ctx)
}
