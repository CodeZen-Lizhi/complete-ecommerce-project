package main

import (
	"context"
	"testing"

	"github.com/cloudwego/eino/schema"
)

func TestGenerateRequiresModel(t *testing.T) {
	t.Parallel()

	if _, err := generate(context.Background(), nil, nil); err == nil {
		t.Fatal("模型为空时应返回错误")
	}
}

func TestAssistantText(t *testing.T) {
	t.Parallel()

	message := &schema.AgenticMessage{
		Role: schema.AgenticRoleTypeAssistant,
		ContentBlocks: []*schema.ContentBlock{
			schema.NewContentBlock(&schema.AssistantGenText{Text: "第一段"}),
			schema.NewContentBlock(&schema.Reasoning{Text: "不应输出"}),
			schema.NewContentBlock(&schema.AssistantGenText{Text: "第二段"}),
		},
	}

	text, err := assistantText(message)
	if err != nil {
		t.Fatalf("assistantText 返回错误: %v", err)
	}
	if text != "第一段\n第二段" {
		t.Fatalf("text = %q", text)
	}
}

func TestAssistantTextWithoutGeneratedText(t *testing.T) {
	t.Parallel()

	_, err := assistantText(&schema.AgenticMessage{})
	if err == nil {
		t.Fatal("缺少 AssistantGenText 时应返回错误")
	}
}
