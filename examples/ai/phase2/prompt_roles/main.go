package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cloudwego/eino-ext/components/model/agenticopenai"
	einomodel "github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

const (
	baseURL        = "http://localhost:8084/v1"
	modelName      = "gpt-5.4-mini"
	apiKey         = "replace-with-your-api-key"
	requestTimeout = 30 * time.Second
)

type promptScenario struct {
	name         string
	systemPrompt string
}

var _ einomodel.AgenticModel = (*agenticopenai.ResponsesModel)(nil)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	scenarios := []promptScenario{
		{
			name:         "初学者导师",
			systemPrompt: "你是一名面向 Go 初学者的导师，请使用通俗语言、类比和一个短小代码示例。",
		},
		{
			name:         "资深工程师顾问",
			systemPrompt: "你是一名面向资深 Go 工程师的技术顾问，请强调接口边界、抽象成本和适用条件。",
		},
	}
	const userPrompt = "请解释 Go 的 interface 适合解决什么问题，并给出一个简短例子。"

	startedAt := time.Now()
	agenticModel, err := newAgenticModel(ctx)
	if err != nil {
		fmt.Printf("创建 AgenticModel 失败: %v\n", err)
		return
	}

	for _, scenario := range scenarios {
		messages := []*schema.AgenticMessage{
			schema.SystemAgenticMessage(scenario.systemPrompt),
			schema.UserAgenticMessage(userPrompt),
		}

		response, err := generate(ctx, agenticModel, messages)
		if err != nil {
			fmt.Printf("场景 %q 模型调用失败: %v\n", scenario.name, err)
			return
		}

		answer, err := assistantText(response)
		if err != nil {
			fmt.Printf("场景 %q 读取模型回答失败: %v\n", scenario.name, err)
			return
		}

		fmt.Printf("场景：%s\n", scenario.name)
		fmt.Printf("System Prompt：%s\n", scenario.systemPrompt)
		fmt.Printf("模型回答：\n%s\n\n", answer)
	}

	fmt.Printf("总耗时: %s\n", time.Since(startedAt))
}

func newAgenticModel(ctx context.Context) (einomodel.AgenticModel, error) {
	if strings.TrimSpace(apiKey) == "" || apiKey == "replace-with-your-api-key" {
		return nil, fmt.Errorf("API Key 未配置")
	}

	store := false
	maxRetries := 0
	timeout := requestTimeout
	config := agenticopenai.ResponsesConfig{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		Model:      modelName,
		Timeout:    &timeout,
		MaxRetries: &maxRetries,
		Store:      &store,
	}

	agenticModel, err := agenticopenai.NewResponsesModel(ctx, &config)
	if err != nil {
		return nil, fmt.Errorf("创建 Eino AgenticModel 失败: %w", err)
	}

	return agenticModel, nil
}

func generate(
	ctx context.Context,
	agenticModel einomodel.AgenticModel,
	messages []*schema.AgenticMessage,
) (*schema.AgenticMessage, error) {
	if agenticModel == nil {
		return nil, fmt.Errorf("AgenticModel 不能为空")
	}
	if len(messages) == 0 {
		return nil, fmt.Errorf("消息列表不能为空")
	}

	response, err := agenticModel.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("Eino Generate 调用失败: %w", err)
	}
	if response == nil {
		return nil, fmt.Errorf("响应为空")
	}

	return response, nil
}

func assistantText(message *schema.AgenticMessage) (string, error) {
	if message == nil {
		return "", fmt.Errorf("AgenticMessage 不能为空")
	}

	parts := make([]string, 0, len(message.ContentBlocks))
	for _, block := range message.ContentBlocks {
		if block == nil || block.Type != schema.ContentBlockTypeAssistantGenText || block.AssistantGenText == nil {
			continue
		}
		if text := strings.TrimSpace(block.AssistantGenText.Text); text != "" {
			parts = append(parts, text)
		}
	}
	if len(parts) == 0 {
		return "", fmt.Errorf("响应中没有 AssistantGenText")
	}

	return strings.Join(parts, "\n"), nil
}
