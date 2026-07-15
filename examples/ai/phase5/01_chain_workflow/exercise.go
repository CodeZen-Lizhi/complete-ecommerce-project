package main

import (
	"context"
	"errors"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type chainInput struct {
	Question string
}

type chainOutput struct {
	Answer string
}

type chainNode[I any, O any] interface {
	// Invoke 执行一个确定性节点并传播 Context。
	Invoke(ctx context.Context, input I) (O, error)
}

type compiledChain interface {
	// Invoke 执行已经编译的固定顺序工作流。
	Invoke(ctx context.Context, input chainInput) (chainOutput, error)
}

// buildChain 按规范化、Prompt、模型、解析的顺序构建并编译 Chain。
func buildChain(ctx context.Context) (compiledChain, error) {
	return nil, errExerciseIncomplete
}

// runExercise 按执行顺序组织“固定顺序 Chain”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：定义每个节点的输入输出类型和唯一职责。
	// TODO 2：使用 Eino Chain 按规范化、Prompt、模型调用和解析顺序组装节点。
	// TODO 3：编译 Chain，并让节点错误和 Context 取消完整传播。
	// TODO 4：使用固定输入调用 compiledChain.Invoke。
	// TODO 5：记录节点耗时，覆盖空输入和任一中间节点失败。
	return errExerciseIncomplete
}
