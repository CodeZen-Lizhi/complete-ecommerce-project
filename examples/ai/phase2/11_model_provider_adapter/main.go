package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type message struct {
	Role    string
	Content string
}

type chatResult struct {
	Content      string
	InputTokens  int
	OutputTokens int
}

type chatProvider interface {
	Generate(ctx context.Context, messages []message) (chatResult, error)
}

type providerConfig struct {
	Provider string
	BaseURL  string
	APIKey   string
	Model    string
}

// main 演示应用只依赖 chatProvider，而不依赖具体模型 SDK。
func main() {
	config := providerConfig{
		Provider: "openai-compatible",
		BaseURL:  "http://localhost:8084/v1",
		APIKey:   "replace-with-your-api-key",
		Model:    "gpt-5.4-mini",
	}

	ctx := context.Background()
	provider, err := newProvider(ctx, config)
	if err != nil {
		fmt.Printf("创建 Provider 失败: %v\n", err)
		return
	}

	result, err := runDemo(ctx, provider)
	if err != nil {
		fmt.Printf("运行 Provider 适配练习失败: %v\n", err)
		return
	}
	fmt.Printf("模型回答: %s\n输入 Token: %d\n输出 Token: %d\n", result.Content, result.InputTokens, result.OutputTokens)
}

// newProvider 校验配置，并按 provider 名称创建具体模型适配器。
func newProvider(ctx context.Context, config providerConfig) (chatProvider, error) {
	if ctx == nil {
		return nil, errors.New("Context 不能为空")
	}
	if strings.TrimSpace(config.Provider) == "" {
		return nil, errors.New("Provider 不能为空")
	}
	if strings.TrimSpace(config.APIKey) == "" || config.APIKey == "replace-with-your-api-key" {
		return nil, errors.New("API Key 未配置")
	}
	if strings.TrimSpace(config.Model) == "" {
		return nil, errors.New("Model 不能为空")
	}

	// TODO 1：为 OpenAI-compatible Eino ChatModel 定义 adapter，实现 chatProvider。
	// TODO 2：在 adapter 内完成项目 Message 与 schema.Message 的双向转换，错误用 %w 包装。
	// TODO 3：根据 config.Provider 白名单选择 adapter；未知名称明确失败，不做静默 fallback。
	// TODO 4：把 BaseURL、APIKey、Model 和超时传给具体 SDK，但不让业务层读取 SDK 配置。
	// TODO 5：使用 Fake Provider 测试 Factory，并验证切换配置不需要修改调用方。
	return nil, errExerciseIncomplete
}

// runDemo 只通过项目 chatProvider 接口完成一次调用，不依赖具体模型 SDK。
func runDemo(ctx context.Context, provider chatProvider) (chatResult, error) {
	if ctx == nil {
		return chatResult{}, errors.New("Context 不能为空")
	}
	if provider == nil {
		return chatResult{}, errors.New("Chat Provider 不能为空")
	}

	// TODO 6：组装项目 message，调用 provider.Generate，并使用 %w 包装厂商适配器错误。
	// 成功后校验 Content 非空和 Token 非负，再返回统一 chatResult。
	return chatResult{}, errExerciseIncomplete
}
