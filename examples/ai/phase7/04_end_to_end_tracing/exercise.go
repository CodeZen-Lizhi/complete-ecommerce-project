package main

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
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

// newOpenTelemetryProvider 创建真实 stdout Exporter 和 SDK TracerProvider，调用方负责 Shutdown。
func newOpenTelemetryProvider(ctx context.Context) (*sdktrace.TracerProvider, *stdouttrace.Exporter, error) {
	// TODO 1：注册 stdout Exporter 和 SDK TracerProvider，并在入口创建根 Span。
	return nil, nil, errExerciseIncomplete
}

// tracePipeline 为模型、Retriever、Rerank、Tool 和 Agent 创建父子 Span。
func tracePipeline(ctx context.Context, tracer tracer) error {
	// TODO 2：为模型、Retriever、Rerank、Tool 和 Agent 创建正确的父子 Span。
	return errExerciseIncomplete
}

// safeSpanAttributes 生成低敏、低基数的 Span 属性。
func safeSpanAttributes(modelName string, resultCount int) (spanAttributes, error) {
	// TODO 3：只记录耗时、状态、模型名和结果数量等低敏属性。
	return nil, errExerciseIncomplete
}

// recordSpanFailure 写入错误状态但不泄露敏感内容。
func recordSpanFailure(target span, cause error) error {
	// TODO 4：记录错误状态，但不写入 Prompt、Secret 或完整文档。
	return errExerciseIncomplete
}

// verifyTraceTree 验证一次请求可以还原完整调用链。
func verifyTraceTree(ctx context.Context) error {
	// TODO 5：执行管线并检查父子关系、错误状态和 Provider Shutdown。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“端到端 Trace”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return verifyTraceTree(ctx)
}
