package main

import (
	"context"
	"errors"
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

	systemPrompt     = "你是一名 Go 学习助手。回答需要准确、简洁，并延续当前对话中的上下文。"
	firstUserPrompt  = "请用一个生活中的类比解释 Go interface，并给出一个简短代码例子。"
	secondUserPrompt = "把你刚才例子中的 interface 改成一个更贴近电商支付场景的命名，其他结构尽量不变。"
)

var _ einomodel.AgenticModel = (*agenticopenai.ResponsesModel)(nil)

func main() {
	ctx := context.Background()

	agenticModel, err := newAgenticModel(ctx)
	if err != nil {
		fmt.Printf("创建 AgenticModel 失败: %v\n", err)
		return
	}

	if err := runConversation(ctx, agenticModel); err != nil {
		fmt.Printf("运行多轮对话失败: %v\n", err)
	}
}

func newAgenticModel(ctx context.Context) (einomodel.AgenticModel, error) {
	if err := validateAPIKey(apiKey); err != nil {
		return nil, err
	}

	store := false
	maxRetries := 0
	timeout := requestTimeout
	config := agenticopenai.ResponsesConfig{
		BaseURL:         baseURL,
		APIKey:          apiKey,
		Model:           modelName,
		Timeout:         &timeout,
		MaxRetries:      &maxRetries,
		Store:           &store,
		EnableAutoCache: false,
	}

	agenticModel, err := agenticopenai.NewResponsesModel(ctx, &config)
	if err != nil {
		return nil, fmt.Errorf("创建 Eino AgenticModel 失败: %w", err)
	}

	return agenticModel, nil
}

func runConversation(ctx context.Context, agenticModel einomodel.AgenticModel) error {
	if agenticModel == nil {
		return errors.New("AgenticModel 不能为空")
	}
	// TODO 1：初始化消息历史。
	// 创建 []*schema.AgenticMessage，只放入一条使用 systemPrompt 构造的 System Message。
	// 这条 System Message 只添加一次，后续每轮都复用同一个 messages 切片。
	messages := []*schema.AgenticMessage{
		schema.SystemAgenticMessage(systemPrompt),
	}
	// TODO 2：组装第一轮 User Message。
	// 使用 firstUserPrompt 创建 User Message，并通过 append 追加到 messages 末尾。
	// 追加后顺序应为 System -> User 1；不要重新创建一个丢失 System Message 的切片。
	messages = append(messages, schema.UserAgenticMessage(firstUserPrompt))
	// TODO 3：发起第一轮模型调用。
	// 调用 generateTurn，传入 ctx、agenticModel、当前完整 messages 和轮次 1。
	// 检查返回的 err；失败时直接返回，让 generateTurn 包装的轮次和底层错误继续向上传播。
	result1, err := generateTurn(ctx, agenticModel, messages, 1)
	// TODO 4：保存并展示第一轮 Assistant Message。
	// 先调用 assistantText 提取第一轮可展示文本并检查错误，再打印第一轮 User Prompt 和 Assistant 回答。
	// 然后把第一轮返回的完整 Assistant AgenticMessage 追加到 messages，不能只保存提取出的字符串。
	if err != nil {
		return fmt.Errorf("第一轮调用失败: %w", err)
	}
	firstAnswer, err := assistantText(result1)
	if err != nil {
		return fmt.Errorf("第一轮解析失败: %w", err)
	}
	fmt.Printf("第一轮 User Prompt: %s\n", firstUserPrompt)
	fmt.Printf("第一轮 Assistant 回答:\n%s\n\n", firstAnswer)
	messages = append(messages, result1)
	// TODO 5：组装第二轮追问。
	// 使用 secondUserPrompt 创建新的 User Message，并追加到同一个 messages。
	// 此时顺序必须是 System -> User 1 -> Assistant 1 -> User 2；可在调用前检查长度和角色帮助排错。
	messages = append(messages, schema.UserAgenticMessage(secondUserPrompt))
	// TODO 6：发起第二轮模型调用。
	// 再次调用 generateTurn，传入包含全部四条历史消息的 messages 和轮次 2。
	// 不要只传第二轮 User Message，否则模型看不到第一轮回答，也就无法理解“刚才例子”的指代。
	result2, err := generateTurn(ctx, agenticModel, messages, 2)
	// TODO 7：展示第二轮结果并结束。
	// 使用 assistantText 提取第二轮回答并检查错误，再打印第二轮 User Prompt 和 Assistant 回答。
	// 所有步骤成功后返回 nil。
	if err != nil {
		return fmt.Errorf("第二轮调用失败: %w", err)
	}
	secondAnswer, err := assistantText(result2)
	if err != nil {
		return fmt.Errorf("第二轮解析失败: %w", err)
	}
	fmt.Printf("第二轮 User Prompt: %s\n", secondUserPrompt)
	fmt.Printf("第二轮 Assistant 回答:\n%s\n", secondAnswer)
	return nil
}

func generateTurn(
	ctx context.Context,
	agenticModel einomodel.AgenticModel,
	messages []*schema.AgenticMessage,
	turn int,
) (*schema.AgenticMessage, error) {
	if ctx == nil {
		return nil, errors.New("Context 不能为空")
	}
	if agenticModel == nil {
		return nil, errors.New("AgenticModel 不能为空")
	}
	if len(messages) == 0 {
		return nil, errors.New("消息列表不能为空")
	}
	if turn <= 0 {
		return nil, errors.New("轮次必须大于 0")
	}

	turnCtx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	response, err := agenticModel.Generate(turnCtx, messages)
	if err != nil {
		return nil, fmt.Errorf("第 %d 轮 Eino Generate 调用失败: %w", turn, err)
	}
	if response == nil {
		return nil, fmt.Errorf("第 %d 轮响应为空", turn)
	}

	return response, nil
}

func assistantText(message *schema.AgenticMessage) (string, error) {
	if message == nil {
		return "", errors.New("AgenticMessage 不能为空")
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
		return "", errors.New("响应中没有 AssistantGenText")
	}

	return strings.Join(parts, "\n"), nil
}

func validateAPIKey(key string) error {
	trimmed := strings.TrimSpace(key)
	if trimmed == "" || trimmed == "replace-with-your-api-key" {
		return errors.New("API Key 未配置")
	}
	return nil
}
