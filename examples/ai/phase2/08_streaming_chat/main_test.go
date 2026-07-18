package main

import (
	"context"
	"errors"
	"strings"
	"testing"

	einomodel "github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

type fakeStreamingModel struct {
	reader    *schema.StreamReader[*schema.Message]
	streamErr error
}

// Generate 满足 BaseChatModel 接口；流式消费测试不会调用该方法。
func (m *fakeStreamingModel) Generate(context.Context, []*schema.Message, ...einomodel.Option) (*schema.Message, error) {
	return nil, nil
}

// Stream 返回测试预置的响应流或创建错误。
func (m *fakeStreamingModel) Stream(context.Context, []*schema.Message, ...einomodel.Option) (*schema.StreamReader[*schema.Message], error) {
	return m.reader, m.streamErr
}

// TestStreamAnswer 验证正常 EOF 后成功输出全部非空文本块。
func TestStreamAnswer(t *testing.T) {
	t.Parallel()

	model := &fakeStreamingModel{reader: schema.StreamReaderFromArray([]*schema.Message{
		schema.AssistantMessage("你", nil),
		schema.AssistantMessage("", nil),
		schema.AssistantMessage("好", nil),
	})}
	var output strings.Builder

	err := streamAnswer(context.Background(), model, []*schema.Message{schema.UserMessage("问候")}, &output)
	if err != nil {
		t.Fatalf("正常流返回错误: %v", err)
	}
	if output.String() != "你好" {
		t.Fatalf("输出 = %q，期望 %q", output.String(), "你好")
	}
}

// TestStreamAnswerRejectsEmptyOrNilStream 验证空流和 nil chunk 都会明确失败。
func TestStreamAnswerRejectsEmptyOrNilStream(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		chunks      []*schema.Message
		wantErrText string
	}{
		{name: "空流", chunks: nil, wantErrText: "流式响应为空"},
		{name: "nil chunk", chunks: []*schema.Message{nil}, wantErrText: "流式响应块不能为空"},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			model := &fakeStreamingModel{reader: schema.StreamReaderFromArray(test.chunks)}
			err := streamAnswer(context.Background(), model, []*schema.Message{schema.UserMessage("问题")}, &strings.Builder{})
			if err == nil || !strings.Contains(err.Error(), test.wantErrText) {
				t.Fatalf("期望错误包含 %q，实际为: %v", test.wantErrText, err)
			}
		})
	}
}

// TestStreamAnswerPreservesStreamError 验证创建流失败时保留底层错误链。
func TestStreamAnswerPreservesStreamError(t *testing.T) {
	t.Parallel()

	streamErr := errors.New("provider unavailable")
	model := &fakeStreamingModel{streamErr: streamErr}

	err := streamAnswer(context.Background(), model, []*schema.Message{schema.UserMessage("问题")}, &strings.Builder{})
	if !errors.Is(err, streamErr) || !strings.Contains(err.Error(), "创建流式响应失败") {
		t.Fatalf("创建流错误没有正确包装: %v", err)
	}
}
