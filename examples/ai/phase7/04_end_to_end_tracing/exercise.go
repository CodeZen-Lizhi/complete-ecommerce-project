package main

import (
	"context"
	"errors"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type traceContext struct {
	TraceID string
	SpanID  string
}

type spanAttributes map[string]string

type span interface {
	// End 结束 Span，并记录最终错误状态。
	End(err error)
	// SetAttributes 写入低基数、非敏感属性。
	SetAttributes(attributes spanAttributes)
}

type tracer interface {
	// Start 基于父 Context 创建子 Span。
	Start(ctx context.Context, name string) (context.Context, span)
}

// tracePipeline 为模型、Retriever、Rerank、Tool 和 Agent 创建父子 Span。
func tracePipeline(ctx context.Context, tracer tracer) error {
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“端到端 Trace”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：在入口创建或校验 Trace ID，并写入 Context。
	// TODO 2：通过 tracer 为模型、Retriever、Rerank、Tool 和 Agent 创建父子 Span。
	// TODO 3：只记录耗时、状态、模型名和结果数量等低敏属性。
	// TODO 4：错误写入 Span 状态，但不记录 Prompt、Secret 或完整文档。
	// TODO 5：实现 tracePipeline 测试，验证一次请求能还原完整调用链。
	return errExerciseIncomplete
}
