package main

import (
	"context"
	"errors"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type callFingerprint struct {
	ToolName            string
	NormalizedArguments string
}

type loopState struct {
	Steps          int
	MaxSteps       int
	RepeatedCalls  map[callFingerprint]int
	NoProgressRuns int
}

// normalizeArguments 规范化 JSON 参数以生成稳定调用指纹。
func normalizeArguments(argumentsJSON string) (string, error) {
	// TODO 1：稳定排序 JSON 对象字段，并拒绝非法 JSON 和尾随内容。
	return "", errExerciseIncomplete
}

// buildCallFingerprint 使用工具名和规范化参数生成稳定指纹。
func buildCallFingerprint(toolName string, normalizedArguments string) (callFingerprint, error) {
	// TODO 2：拒绝空工具名和空参数，并返回可比较指纹。
	return callFingerprint{}, errExerciseIncomplete
}

// recordCall 记录调用并判断是否达到重复或无进展阈值。
func recordCall(state loopState, fingerprint callFingerprint, progressed bool, repeatLimit int) (loopState, error) {
	// TODO 3：统计相同调用和连续无进展次数，不修改传入状态。
	return loopState{}, errExerciseIncomplete
}

// checkLoopTermination 判断步骤、重复和无进展阈值。
func checkLoopTermination(state loopState, repeatLimit int, noProgressLimit int) error {
	// TODO 4：达到任一阈值时返回包含终止原因的明确错误。
	return errExerciseIncomplete
}

// verifyLoopDetection 使用确定性循环验证有限步终止。
func verifyLoopDetection(ctx context.Context) error {
	// TODO 5：构造无限工具循环，并确认在阈值内停止。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“循环与重复调用检测”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return verifyLoopDetection(ctx)
}
