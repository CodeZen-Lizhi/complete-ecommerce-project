package main

import (
	"context"
	"errors"

	"github.com/cloudwego/eino/compose"
)

// 模型配置集中放在顶部，练习时直接替换占位值。
const (
	baseURL   = "http://localhost:8084/v1"
	apiKey    = "replace-with-your-api-key"
	modelName = "gpt-5.4-mini"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

// TODO 1：补全每个节点的输入输出类型，并保持单一职责。
type chainInput struct {
	Question string
}

type chainModelConfig struct {
	BaseURL string
	APIKey  string
	Model   string
}

type chainOutput struct {
	Answer string
}

type chainNode[I any, O any] interface {
	// Invoke 执行一个确定性节点并传播 Context。
	Invoke(ctx context.Context, input I) (O, error)
}

// newChainBuilder 按规范化、Prompt、模型、解析的顺序注册节点。
func newChainBuilder() (*compose.Chain[chainInput, chainOutput], error) {
	// TODO 2：使用 compose.NewChain 注册规范化、Prompt、模型调用和解析节点。
	return nil, errExerciseIncomplete
}

// compileChain 编译 Chain 并保留错误与 Context 取消语义。
func compileChain(ctx context.Context, builder *compose.Chain[chainInput, chainOutput]) (compose.Runnable[chainInput, chainOutput], error) {
	// TODO 3：调用 Compile，并让节点错误与 Context 取消完整传播。
	return nil, errExerciseIncomplete
}

// buildChain 构建并编译完整 Chain。
func buildChain(ctx context.Context) (compose.Runnable[chainInput, chainOutput], error) {
	builder, err := newChainBuilder()
	if err != nil {
		return nil, err
	}
	return compileChain(ctx, builder)
}

// invokeChain 使用固定输入调用编译后的 Chain。
func invokeChain(ctx context.Context, chain compose.Runnable[chainInput, chainOutput], input chainInput) (chainOutput, error) {
	// TODO 4：调用 chain.Invoke，并校验空输入和空输出。
	return chainOutput{}, errExerciseIncomplete
}

// observeChainRun 记录节点耗时并验证中间节点失败路径。
func observeChainRun(ctx context.Context, config chainModelConfig) error {
	// TODO 5：输出各节点耗时，覆盖空输入和任一中间节点失败。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“固定顺序 Chain”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return observeChainRun(ctx, chainModelConfig{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Model:   modelName,
	})
}
