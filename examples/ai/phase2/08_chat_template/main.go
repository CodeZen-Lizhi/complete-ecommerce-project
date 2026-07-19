package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

type templateInput struct {
	Role     string
	Question string
	History  []*schema.Message
}

// main 运行 ChatTemplate 变量渲染练习。
func main() {
	input := templateInput{
		Role:     "Go 学习助手",
		Question: "请解释 context.Context 的取消传播。",
		History:  nil,
	}

	messages, err := renderMessages(context.Background(), input)
	if err != nil {
		fmt.Printf("渲染 Prompt 失败: %v\n", err)
		return
	}

	for _, message := range messages {
		fmt.Printf("%s: %s\n", message.Role, message.Content)
	}
}

// buildTemplate 创建包含 System、历史占位符和 User 变量的 Eino ChatTemplate。
func buildTemplate() (prompt.ChatTemplate, error) {
	// TODO 1：调用 prompt.FromMessages(schema.FString, ...) 创建模板。
	// System Message 使用 {role}，中间加入 schema.MessagesPlaceholder("history", true)，
	// User Message 使用 {question}；第二个参数 true 表示历史为空时允许跳过。
	chatTemplate := prompt.FromMessages(schema.FString,
		schema.SystemMessage("你是{role}"),
		schema.MessagesPlaceholder("history", true),
		schema.UserMessage("{question}"),
	)
	return chatTemplate, nil
}

// renderMessages 校验变量并使用 ChatTemplate 生成最终消息列表。
func renderMessages(ctx context.Context, input templateInput) ([]*schema.Message, error) {
	if ctx == nil {
		return nil, errors.New("Context 不能为空")
	}
	if strings.TrimSpace(input.Role) == "" {
		return nil, errors.New("角色不能为空")
	}
	if strings.TrimSpace(input.Question) == "" {
		return nil, errors.New("问题不能为空")
	}

	// TODO 2：调用 buildTemplate 并检查错误。
	template, err := buildTemplate()
	if err != nil {
		return nil, fmt.Errorf("创建 ChatTemplate 失败: %w", err)
	}
	variables := map[string]any{
		"role":     input.Role,
		"history":  input.History,
		"question": input.Question,
	}

	// TODO 3：构造 map[string]any，键必须与模板中的 role、history、question 完全一致。
	// TODO 4：调用 template.Format(ctx, variables)，用 %w 包装缺少变量或格式化失败。
	messages, err := template.Format(ctx, variables)
	if err != nil {
		return nil, fmt.Errorf("格式化 ChatTemplate 失败: %w", err)
	}
	if len(messages) == 0 {
		return nil, errors.New("格式化结果不能为空")
	}
	// TODO 5：检查结果顺序为 System -> 可选历史 -> User，并返回 messages。
	if len(messages) != len(input.History)+2 {
		return nil, fmt.Errorf("消息数量为 %d，期望 %d", len(messages), len(input.History)+2)
	}
	if messages[0] == nil || messages[0].Role != schema.System {
		return nil, errors.New("第一条消息必须是 System Message")
	}
	for index, historyMessage := range input.History {
		if historyMessage == nil || messages[index+1] == nil {
			return nil, fmt.Errorf("第 %d 条历史消息不能为空", index)
		}
		if messages[index+1].Role != historyMessage.Role {
			return nil, fmt.Errorf("第 %d 条历史消息角色顺序不一致", index)
		}
	}
	if messages[len(messages)-1] == nil || messages[len(messages)-1].Role != schema.User {
		return nil, errors.New("最后一条消息必须是 User Message")
	}
	return messages, nil
}
