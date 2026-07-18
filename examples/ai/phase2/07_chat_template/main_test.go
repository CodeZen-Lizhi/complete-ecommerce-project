package main

import (
	"context"
	"testing"

	"github.com/cloudwego/eino/schema"
)

// TestRenderMessages 验证模板变量渲染和消息顺序。
func TestRenderMessages(t *testing.T) {
	t.Parallel()

	history := []*schema.Message{
		schema.UserMessage("历史问题"),
		schema.AssistantMessage("历史回答", nil),
	}
	messages, err := renderMessages(context.Background(), templateInput{
		Role:     "Go 学习助手",
		Question: "什么是 Context？",
		History:  history,
	})
	if err != nil {
		t.Fatalf("renderMessages 返回错误: %v", err)
	}
	if len(messages) != 4 {
		t.Fatalf("消息数量 = %d，期望 4", len(messages))
	}
	if messages[0].Role != schema.System || messages[0].Content != "你是Go 学习助手" {
		t.Fatalf("System Message 渲染错误: %#v", messages[0])
	}
	if messages[1] != history[0] || messages[2] != history[1] {
		t.Fatalf("历史消息没有保持原有顺序")
	}
	if messages[3].Role != schema.User || messages[3].Content != "什么是 Context？" {
		t.Fatalf("User Message 渲染错误: %#v", messages[3])
	}
}

// TestRenderMessagesWithoutHistory 验证可选历史为空时仍生成 System 和 User 消息。
func TestRenderMessagesWithoutHistory(t *testing.T) {
	t.Parallel()

	messages, err := renderMessages(context.Background(), templateInput{
		Role:     "助手",
		Question: "问题",
	})
	if err != nil {
		t.Fatalf("renderMessages 返回错误: %v", err)
	}
	if len(messages) != 2 || messages[0].Role != schema.System || messages[1].Role != schema.User {
		t.Fatalf("空历史消息顺序错误: %#v", messages)
	}
}
