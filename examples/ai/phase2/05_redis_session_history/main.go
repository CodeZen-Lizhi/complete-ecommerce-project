package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cloudwego/eino-ext/components/model/agenticopenai"
	einomodel "github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/redis/go-redis/v9"
)

const (
	baseURL         = "http://localhost:8084/v1"
	modelName       = "gpt-5.4-mini"
	apiKey          = "replace-with-your-api-key"
	redisAddr       = "localhost:6379"
	requestTimeout  = 30 * time.Second
	exerciseTimeout = 90 * time.Second

	sessionKeyPrefix = "ai:phase2:session:"
	demoSessionID    = "123456789"
	systemPrompt     = "你是一名 Go 学习助手。回答需要准确、简洁，并延续当前会话中的上下文。"
	firstUserPrompt  = "请用一个生活中的类比解释 Go interface，并给出一个简短代码例子。"
	secondUserPrompt = "把你刚才例子中的 interface 改成更贴近电商支付场景的命名，" +
		"其他结构尽量不变。"
)

var _ einomodel.AgenticModel = (*agenticopenai.ResponsesModel)(nil)

type historyStore interface {
	// Load 按 session ID 读取并反序列化完整消息历史；空历史不视为错误。
	Load(ctx context.Context, sessionID string) ([]*schema.AgenticMessage, error)
	// Append 按原顺序追加指定 session 本轮新增的消息。
	Append(ctx context.Context, sessionID string, messages ...*schema.AgenticMessage) error
}

type redisHistoryStore struct {
	client redis.Cmdable
}

// main 初始化模型与 Redis Client，并运行两轮 Redis 会话历史练习。
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), exerciseTimeout)
	defer cancel()

	agenticModel, err := newAgenticModel(ctx)
	if err != nil {
		fmt.Printf("创建 AgenticModel 失败: %v\n", err)
		return
	}

	redisClient := redis.NewClient(&redis.Options{Addr: redisAddr})
	defer func() {
		if closeErr := redisClient.Close(); closeErr != nil {
			fmt.Printf("关闭 Redis Client 失败: %v\n", closeErr)
		}
	}()

	store, err := newRedisHistoryStore(redisClient)
	if err != nil {
		fmt.Printf("创建 Redis 历史存储失败: %v\n", err)
		return
	}

	if err := runExercise(ctx, agenticModel, store); err != nil {
		fmt.Printf("运行 Redis 会话历史练习失败: %v\n", err)
	}
}

// newAgenticModel 创建关闭服务端存储和自动缓存的 Eino AgenticModel。
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

// newRedisHistoryStore 使用给定 Redis Client 创建会话历史存储。
func newRedisHistoryStore(client redis.Cmdable) (*redisHistoryStore, error) {
	if client == nil {
		return nil, errors.New("Redis Client 不能为空")
	}
	return &redisHistoryStore{client: client}, nil
}

// sessionKey 校验 session ID，并生成带固定命名空间的 Redis Key。
func sessionKey(sessionID string) (string, error) {
	// TODO 1：校验 session ID 并生成 Redis Key。
	// 去掉 sessionID 首尾空白后，拒绝空字符串；再逐个检查字符，只允许英文字母、数字、短横线和下划线。
	// 校验通过后返回 sessionKeyPrefix + 清洗后的 session ID，避免不同业务的 Key 混在同一命名空间。
	trimmed := strings.TrimSpace(sessionID)
	if trimmed == "" {
		return "", errors.New("session ID 不能为空")
	}
	for _, char := range trimmed {
		if (char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-' || char == '_' {
			continue
		}
		return "", fmt.Errorf("session ID 包含非法字符 %q", char)
	}

	return sessionKeyPrefix + trimmed, nil
}

// Load 按 session ID 读取并反序列化完整消息历史；空历史不视为错误。
func (s *redisHistoryStore) Load(
	ctx context.Context,
	sessionID string,
) ([]*schema.AgenticMessage, error) {
	if ctx == nil {
		return nil, errors.New("Context 不能为空")
	}
	if s == nil || s.client == nil {
		return nil, errors.New("Redis 历史存储未初始化")
	}

	// TODO 2：读取当前 session 的全部原始历史。
	// 先调用 sessionKey 得到 key 并检查错误，再使用 s.client.LRange(ctx, key, 0, -1).Result() 读取整个 List。
	// Redis 调用失败时返回 nil，并用 %w 包装错误；空 List 是合法的新会话，应继续返回空消息切片。
	key, err := sessionKey(sessionID)
	if err != nil {
		return nil, fmt.Errorf("生成 session %q 的 Redis Key 失败: %w", sessionID, err)
	}
	result, err := s.client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("读取 session %q 的 Redis 历史失败: %w", sessionID, err)
	}
	if len(result) == 0 {
		return nil, nil
	}
	// TODO 3：按 Redis 返回顺序反序列化消息。
	// 先在 import 中加入标准库 encoding/json。创建容量等于原始元素数量的 []*schema.AgenticMessage；逐条声明非 nil 指针，
	// 使用 json.Unmarshal 解析每个元素并检查错误，然后依次 append。错误中要包含 session ID 和元素下标。
	messages := make([]*schema.AgenticMessage, 0, len(result))
	for index, value := range result {
		message := new(schema.AgenticMessage)
		if err := json.Unmarshal([]byte(value), message); err != nil {
			return nil, fmt.Errorf(
				"反序列化 session %q 的第 %d 条消息失败: %w",
				sessionID,
				index,
				err,
			)
		}
		messages = append(messages, message)
	}
	return messages, nil
}

// Append 按原顺序序列化并追加指定 session 本轮新增的消息。
func (s *redisHistoryStore) Append(
	ctx context.Context,
	sessionID string,
	messages ...*schema.AgenticMessage,
) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}
	if s == nil || s.client == nil {
		return errors.New("Redis 历史存储未初始化")
	}
	if len(messages) == 0 {
		return errors.New("待保存消息不能为空")
	}

	// TODO 4：校验并序列化本轮新增消息。
	// 先调用 sessionKey 得到 key；创建 []any 保存 RPUSH 参数。按顺序遍历 messages，拒绝 nil 消息，
	// 使用 json.Marshal 编码完整 AgenticMessage，并把每段 JSON 追加到参数切片。编码错误要包含消息下标。
	key, err := sessionKey(sessionID)
	if err != nil {
		return fmt.Errorf("生成 session %q 的 Redis Key 失败: %w", sessionID, err)
	}
	values := make([]any, 0, len(messages))
	for index, message := range messages {
		if message == nil {
			return fmt.Errorf("session %q 的第 %d 条待保存消息不能为空", sessionID, index)
		}
		encoded, err := json.Marshal(message)
		if err != nil {
			return fmt.Errorf(
				"序列化 session %q 的第 %d 条消息失败: %w",
				sessionID,
				index,
				err,
			)
		}
		values = append(values, encoded)
	}
	// TODO 5：一次性追加本轮消息。
	// 调用 s.client.RPush(ctx, key, values...).Err()，保证参数中的 System/User/Assistant 顺序原样进入 List。
	// Redis 写入失败时用 %w 包装并带上 session ID；成功时返回 nil。
	if err := s.client.RPush(ctx, key, values...).Err(); err != nil {
		return fmt.Errorf("追加 session %q 的 Redis 历史失败: %w", sessionID, err)
	}
	return nil
}

// runTurn 从 Redis 重组指定 session 的历史，执行一轮模型调用并保存新增消息。
func runTurn(
	ctx context.Context,
	agenticModel einomodel.AgenticModel,
	store historyStore,
	sessionID string,
	userPrompt string,
	turn int,
) (string, error) {
	if ctx == nil {
		return "", errors.New("Context 不能为空")
	}
	if agenticModel == nil {
		return "", errors.New("AgenticModel 不能为空")
	}
	if store == nil {
		return "", errors.New("历史存储不能为空")
	}
	if strings.TrimSpace(userPrompt) == "" {
		return "", errors.New("User Prompt 不能为空")
	}
	if turn <= 0 {
		return "", errors.New("轮次必须大于 0")
	}

	// TODO 6：从 Redis 加载当前 session 的历史。
	// 调用 store.Load(ctx, sessionID) 并检查错误；不要复用上一轮函数调用留下的内存切片。
	messages, err := store.Load(ctx, sessionID)
	if err != nil {
		return "", fmt.Errorf("第 %d 轮加载 session %q 历史失败: %w", turn, sessionID, err)
	}
	// TODO 7：初始化本轮待保存消息。
	pendingMessages := make([]*schema.AgenticMessage, 0, 3)
	// 创建 pendingMessages；如果加载结果为空，向历史和 pendingMessages 各追加同一条 System Message。
	// 如果已有历史，不要重复添加 System Message。
	if len(messages) == 0 {
		systemMessage := schema.SystemAgenticMessage(systemPrompt)
		pendingMessages = append(pendingMessages, systemMessage)
		messages = append(messages, systemMessage)
	}
	userAgenticMessage := schema.UserAgenticMessage(userPrompt)
	pendingMessages = append(pendingMessages, userAgenticMessage)
	messages = append(messages, userAgenticMessage)
	// TODO 8：追加本轮 User 并调用模型。
	// 创建 User Message，同时追加到完整历史和 pendingMessages；再调用 generateTurn，传入包含当前 User 的完整历史。
	// 检查模型错误，失败时直接返回，不要把没有 Assistant 回答的半轮对话写入 Redis。
	result, err := generateTurn(ctx, agenticModel, messages, turn)
	if err != nil {
		return "", err
	}
	pendingMessages = append(pendingMessages, result)
	// TODO 9：校验并保留完整 Assistant 消息。
	// 先调用 assistantText 提取可展示文本并检查错误；成功后再把完整 Assistant AgenticMessage 追加到 pendingMessages。
	// 文本提取失败时不要写 Redis，错误必须带轮次上下文。
	text, err := assistantText(result)
	if err != nil {
		return "", fmt.Errorf("第 %d 轮解析 Assistant 响应失败: %w", turn, err)
	}
	// TODO 10：保存本轮新增消息并返回文本。
	// 只调用一次 store.Append 保存 pendingMessages，不能保存完整历史，否则第二轮会把第一轮消息重复写入 Redis。
	// Redis 保存失败时带轮次上下文并返回错误；成功时返回 TODO 9 提取出的 Assistant 文本。
	err = store.Append(ctx, sessionID, pendingMessages...)
	if err != nil {
		return "", fmt.Errorf("第 %d 轮保存 session %q 历史失败: %w", turn, sessionID, err)
	}
	return text, nil
}

// runExercise 使用同一个 session ID 顺序执行两轮相互关联的对话。
func runExercise(
	ctx context.Context,
	agenticModel einomodel.AgenticModel,
	store historyStore,
) error {
	// TODO 11：执行第一轮独立请求。
	// 调用 runTurn，传入 demoSessionID、firstUserPrompt 和轮次 1；检查错误并打印第一轮 User/Assistant。
	firstAnswer, err := runTurn(ctx, agenticModel, store, demoSessionID, firstUserPrompt, 1)
	if err != nil {
		return fmt.Errorf("第一轮对话失败: %w", err)
	}
	fmt.Printf("第一轮 User Prompt: %s\n", firstUserPrompt)
	fmt.Printf("第一轮 Assistant 回答:\n%s\n\n", firstAnswer)
	// TODO 12：模拟进程内没有历史的第二轮请求。
	// 再次调用 runTurn，仍传入同一个 demoSessionID，但使用 secondUserPrompt 和轮次 2。
	// runTurn 必须重新从 Redis 加载历史；检查错误并打印第二轮 User/Assistant，全部成功后返回 nil。
	result2, err := runTurn(ctx, agenticModel, store, demoSessionID, secondUserPrompt, 2)
	if err != nil {
		return fmt.Errorf("第二轮对话失败: %w", err)
	}
	fmt.Printf("第二轮 User Prompt: %s\n", secondUserPrompt)
	fmt.Printf("第二轮 Assistant 回答:\n%s\n", result2)
	return nil
}

// generateTurn 在独立超时内调用模型，并为错误补充轮次上下文。
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

// assistantText 提取并合并 AgenticMessage 中可展示的 Assistant 文本块。
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

// validateAPIKey 拒绝空值和示例占位符，避免发起无效模型请求。
func validateAPIKey(key string) error {
	trimmed := strings.TrimSpace(key)
	if trimmed == "" || trimmed == "replace-with-your-api-key" {
		return errors.New("API Key 未配置")
	}
	return nil
}
