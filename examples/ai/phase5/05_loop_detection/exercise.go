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
	return "", errExerciseIncomplete
}

// recordCall 记录调用并判断是否达到重复或无进展阈值。
func recordCall(state loopState, fingerprint callFingerprint, progressed bool, repeatLimit int) (loopState, error) {
	return loopState{}, errExerciseIncomplete
}

// runExercise 按执行顺序组织“循环与重复调用检测”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：实现 normalizeArguments，稳定排序 JSON 对象字段并拒绝非法 JSON。
	// TODO 2：为工具名和规范化参数生成 callFingerprint。
	// TODO 3：实现 recordCall，统计相同调用和连续无进展次数。
	// TODO 4：达到最大步骤、重复阈值或无进展阈值时返回明确终止错误。
	// TODO 5：使用确定性 Fake 构造无限循环，验证有限步内停止。
	return errExerciseIncomplete
}
