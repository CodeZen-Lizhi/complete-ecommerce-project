package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cloudwego/eino-ext/components/model/openai"
	einomodel "github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

const (
	baseURL   = "http://localhost:8084/v1"
	apiKey    = "replace-with-your-api-key"
	modelName = "gpt-5.4-mini"
)

var _ einomodel.BaseChatModel = (*openai.ChatModel)(nil)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	messages := []*schema.Message{
		schema.SystemMessage("你是一个 Go 学习助手，请用简洁中文回答。"),
		schema.UserMessage("请用一个例子解释 Go 的 context.Context 有什么作用。"),
	}

	startedAt := time.Now()
	chatModel, err := newChatModel(ctx)
	if err != nil {
		fmt.Printf("创建 ChatModel 失败: %v\n", err)
		return
	}

	response, err := generate(ctx, chatModel, messages)
	if err != nil {
		fmt.Printf("模型调用失败: %v\n", err)
		return
	}

	fmt.Printf("模型回答: %s\n", response.Content)
	fmt.Printf("总耗时: %s\n", time.Since(startedAt))
}

func newChatModel(ctx context.Context) (einomodel.BaseChatModel, error) {
	// TODO 1：检查 apiKey 是否为空或仍是占位符。
	// 配置无效时返回 nil 和一个非 nil error，不能让调用方把空模型误认为创建成功。
	if strings.TrimSpace(apiKey) == "" || apiKey == "replace-with-your-api-key" {
		return nil, fmt.Errorf("API Key 未配置")
	}
	// TODO 2：确认骨架导入的 Eino OpenAI-compatible 模型组件。
	// 包路径是 github.com/cloudwego/eino-ext/components/model/openai；它负责封装底层 HTTP 调用。
	// TODO 3：创建 openai.ChatModelConfig，设置 APIKey、BaseURL、Model 和 Timeout。
	// APIKey 使用 apiKey，BaseURL 使用 baseURL，Model 使用 modelName，Timeout 先设置为 30 秒。
	// 本练习不自定义 HTTPClient；如果设置了 HTTPClient，组件配置中的 Timeout 将不再生效。
	config := openai.ChatModelConfig{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Model:   modelName,
		Timeout: 30 * time.Second,
	}
	// TODO 4：调用 openai.NewChatModel(ctx, config) 创建模型，并检查返回的 err。
	// 创建失败时使用 fmt.Errorf 和 %w 增加“创建 Eino ChatModel 失败”上下文。
	chatModel, err := openai.NewChatModel(ctx, &config)
	if err != nil {
		return nil, fmt.Errorf("创建 Eino ChatModel 失败: %w", err)
	}
	// TODO 5：创建成功后，将模型作为 einomodel.BaseChatModel 返回，error 返回 nil。
	// 本练习只依赖 Generate/Stream 所需的最小接口，不使用已弃用的 einomodel.ChatModel。
	return chatModel, nil
}

func generate(ctx context.Context, chatModel einomodel.BaseChatModel, messages []*schema.Message) (*schema.Message, error) {
	// TODO 6：检查 chatModel 是否为 nil，并检查 messages 是否为空。
	// 任一条件不满足都返回 nil 和明确错误，避免在 SDK 内部发生难以定位的失败。
	if chatModel == nil {
		return nil, fmt.Errorf("ChatModel 不能为空")
	}
	if len(messages) == 0 {
		return nil, fmt.Errorf("消息列表不能为空")
	}
	// TODO 7：调用 chatModel.Generate 发送 messages。
	// 第一个参数传 ctx，第二个参数传 messages。当前模型的 Temperature 固定为 1，因此不传额外选项。
	result, err := chatModel.Generate(ctx, messages)
	// TODO 8：检查 Generate 返回的 err。
	// 调用失败时使用 fmt.Errorf 和 %w 增加“Eino Generate 调用失败”上下文。
	if err != nil {
		return nil, fmt.Errorf("Eino Generate 调用失败: %w", err)
	}
	// TODO 9：检查返回的 response 是否为 nil。
	// SDK 没有返回错误但响应为空时，也必须返回明确错误，不能继续访问 response.Content。
	if result == nil {
		return nil, fmt.Errorf("响应为空")
	}
	// TODO 10：成功后返回 response 和 nil。
	return result, nil
}
