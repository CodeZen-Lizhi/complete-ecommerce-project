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

type failingWriter struct {
	err error
}

// Write 返回测试预置错误，用于验证流式输出失败路径。
func (w *failingWriter) Write([]byte) (int, error) {
	return 0, w.err
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

// TestStreamAnswerPreservesReceiveErrors 验证取消、超时和普通接收错误均保留错误链与分类文案。
func TestStreamAnswerPreservesReceiveErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		streamErr   error
		wantErrText string
	}{
		{name: "取消", streamErr: context.Canceled, wantErrText: "流式请求被取消"},
		{name: "超时", streamErr: context.DeadlineExceeded, wantErrText: "流式请求超时"},
		{name: "中途错误", streamErr: errors.New("连接中断"), wantErrText: "接收流式响应失败"},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			reader, writer := schema.Pipe[*schema.Message](1)
			writer.Send(nil, test.streamErr)
			writer.Close()
			model := &fakeStreamingModel{reader: reader}

			err := streamAnswer(context.Background(), model, []*schema.Message{schema.UserMessage("问题")}, &strings.Builder{})
			if !errors.Is(err, test.streamErr) || !strings.Contains(err.Error(), test.wantErrText) {
				t.Fatalf("接收错误处理不符合预期: %v", err)
			}
		})
	}
}

// TestStreamAnswerPreservesWriterError 验证输出失败时保留 Writer 的底层错误。
func TestStreamAnswerPreservesWriterError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("写入失败")
	model := &fakeStreamingModel{reader: schema.StreamReaderFromArray([]*schema.Message{
		schema.AssistantMessage("回答", nil),
	})}

	err := streamAnswer(
		context.Background(),
		model,
		[]*schema.Message{schema.UserMessage("问题")},
		&failingWriter{err: wantErr},
	)
	if !errors.Is(err, wantErr) || !strings.Contains(err.Error(), "写入流式响应失败") {
		t.Fatalf("Writer 错误处理不符合预期: %v", err)
	}
}
