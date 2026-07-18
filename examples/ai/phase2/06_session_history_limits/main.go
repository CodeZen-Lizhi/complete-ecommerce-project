package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
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
	apiKeyEnv       = "OPENAI_API_KEY"
	redisAddr       = "localhost:6379"
	requestTimeout  = 30 * time.Second
	exerciseTimeout = 90 * time.Second

	maxHistoryTurns  = 2
	sessionTTL       = 30 * time.Minute
	sessionKeyPrefix = "ai:phase2:user-session:"

	demoUserID        = "user-1001"
	otherDemoUserID   = "user-2002"
	demoSessionID     = "checkout-help"
	systemPrompt      = "你是一名 Go 学习助手。回答需要准确、简洁，并延续当前用户会话中的上下文。"
	firstUserPrompt   = "请用一个生活中的类比解释 Go interface。"
	secondUserPrompt  = "把刚才的类比改成电商支付场景，其他结构尽量不变。"
	isolatedUserQuery = "我刚才让你把哪个类比改成电商支付场景？"
)

var (
	_ einomodel.AgenticModel = (*agenticopenai.ResponsesModel)(nil)
)

type historyStore interface {
	// LoadRecent 读取指定用户会话最近的完整 User/Assistant 消息对；空历史不视为错误。
	LoadRecent(ctx context.Context, userID string, sessionID string) ([]*schema.AgenticMessage, error)
	// AppendTurn 原子追加一轮 User/Assistant，裁剪旧消息并刷新会话 TTL。
	AppendTurn(
		ctx context.Context,
		userID string,
		sessionID string,
		userMessage *schema.AgenticMessage,
		assistantMessage *schema.AgenticMessage,
	) error
}

type redisHistoryStore struct {
	client   redis.Cmdable
	maxTurns int
	ttl      time.Duration
}

// main 初始化模型和 Redis Client，并运行历史截断、过期与用户隔离练习。
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

	store, err := newRedisHistoryStore(redisClient, maxHistoryTurns, sessionTTL)
	if err != nil {
		fmt.Printf("创建 Redis 历史存储失败: %v\n", err)
		return
	}

	if err := runExercise(ctx, agenticModel, store); err != nil {
		fmt.Printf("运行会话历史限制练习失败: %v\n", err)
	}
}

// newAgenticModel 创建关闭服务端存储和自动缓存的 Eino AgenticModel。
func newAgenticModel(ctx context.Context) (einomodel.AgenticModel, error) {
	if ctx == nil {
		return nil, errors.New("Context 不能为空")
	}
	apiKey := strings.TrimSpace(os.Getenv(apiKeyEnv))
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

// newRedisHistoryStore 校验历史上限和 TTL，并创建 Redis 会话历史存储。
func newRedisHistoryStore(
	client redis.Cmdable,
	maxTurns int,
	ttl time.Duration,
) (*redisHistoryStore, error) {
	if client == nil {
		return nil, errors.New("Redis Client 不能为空")
	}
	if maxTurns <= 0 {
		return nil, errors.New("历史轮数上限必须大于 0")
	}
	if ttl <= 0 {
		return nil, errors.New("会话 TTL 必须大于 0")
	}
	return &redisHistoryStore{client: client, maxTurns: maxTurns, ttl: ttl}, nil
}

// conversationKey 校验用户和会话标识，并生成同时包含二者的 Redis Key。
func conversationKey(userID string, sessionID string) (string, error) {
	// TODO 1：清洗并校验 userID 与 sessionID，然后生成隔离的 Redis Key。
	// 分别去掉首尾空白，拒绝空值；逐字符限制为英文字母、数字、短横线和下划线。
	// 返回 sessionKeyPrefix + userID + ":" + sessionID，不能只使用 sessionID，否则不同用户会共享历史。
	normalizedUserID, err := normalizeIdentifier("user ID", userID)
	if err != nil {
		return "", err
	}
	normalizedSessionID, err := normalizeIdentifier("session ID", sessionID)
	if err != nil {
		return "", err
	}
	return sessionKeyPrefix + normalizedUserID + ":" + normalizedSessionID, nil
}

// normalizeIdentifier 去除标识首尾空白，并限制为不会破坏 Redis Key 分段的安全字符。
func normalizeIdentifier(name string, value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", fmt.Errorf("%s 不能为空", name)
	}
	for _, char := range trimmed {
		if (char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-' || char == '_' {
			continue
		}
		return "", fmt.Errorf("%s 包含非法字符 %q", name, char)
	}
	return trimmed, nil
}

// LoadRecent 读取并反序列化最近的完整 User/Assistant 消息对。
func (s *redisHistoryStore) LoadRecent(
	ctx context.Context,
	userID string,
	sessionID string,
) ([]*schema.AgenticMessage, error) {
	if ctx == nil {
		return nil, errors.New("Context 不能为空")
	}
	if s == nil || s.client == nil {
		return nil, errors.New("Redis 历史存储未初始化")
	}

	// TODO 2：生成 Key，并只读取末尾 maxTurns * 2 条消息。
	// 使用 conversationKey 校验身份；计算 maxMessages := int64(s.maxTurns * 2)。
	// 调用 LRange(ctx, key, -maxMessages, -1)，Redis 失败时用 %w 包装；空 List 返回空切片和 nil。
	key, err := conversationKey(userID, sessionID)
	if err != nil {
		return nil, fmt.Errorf(
			"生成用户 %q、session %q 的 Redis Key 失败: %w",
			userID,
			sessionID,
			err,
		)
	}
	maxMessages := int64(s.maxTurns * 2)
	values, err := s.client.LRange(ctx, key, -maxMessages, -1).Result()
	if err != nil {
		return nil, fmt.Errorf(
			"读取用户 %q、session %q 的最近历史失败: %w",
			userID,
			sessionID,
			err,
		)
	}
	return decodeHistoryValues(userID, sessionID, values)
}

// decodeHistoryValues 校验完整轮次，并把 Redis JSON 元素解码为非 nil 消息。
func decodeHistoryValues(
	userID string,
	sessionID string,
	values []string,
) ([]*schema.AgenticMessage, error) {
	if len(values) == 0 {
		return []*schema.AgenticMessage{}, nil
	}
	if len(values)%2 != 0 {
		return nil, fmt.Errorf(
			"用户 %q、session %q 的 Redis 历史包含 %d 条消息，不是完整的 User/Assistant 消息对",
			userID,
			sessionID,
			len(values),
		)
	}
	messages := make([]*schema.AgenticMessage, 0, len(values))
	// TODO 3：按 Redis 返回顺序反序列化消息，并校验消息数量为偶数。
	// 使用 encoding/json 将每个元素解析为非 nil *schema.AgenticMessage，错误中包含用户、session 和下标。
	// Redis 只保存完整轮次，因此奇数条消息表示数据损坏，应返回明确错误，不能把半轮历史发给模型。
	for index, value := range values {
		if strings.TrimSpace(value) == "" {
			return nil, fmt.Errorf(
				"用户 %q、session %q 的第 %d 条历史消息为空",
				userID,
				sessionID,
				index,
			)
		}

		var message *schema.AgenticMessage
		if err := json.Unmarshal([]byte(value), &message); err != nil {
			return nil, fmt.Errorf(
				"反序列化用户 %q、session %q 的第 %d 条历史消息失败: %w",
				userID,
				sessionID,
				index,
				err,
			)
		}
		if message == nil {
			return nil, fmt.Errorf(
				"用户 %q、session %q 的第 %d 条历史消息不能为 null",
				userID,
				sessionID,
				index,
			)
		}
		messages = append(messages, message)
	}
	return messages, nil
}

// AppendTurn 原子追加完整一轮消息，裁剪历史并刷新 TTL。
func (s *redisHistoryStore) AppendTurn(
	ctx context.Context,
	userID string,
	sessionID string,
	userMessage *schema.AgenticMessage,
	assistantMessage *schema.AgenticMessage,
) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}
	if s == nil || s.client == nil {
		return errors.New("Redis 历史存储未初始化")
	}

	// TODO 4：校验并序列化完整一轮消息。
	// 调用 conversationKey；拒绝 nil User 或 Assistant；使用 encoding/json 分别编码完整 AgenticMessage。
	// 两条消息必须一起准备完成后才能写 Redis，任何校验或编码失败都不能产生半轮数据。
	key, err := conversationKey(userID, sessionID)
	if err != nil {
		return fmt.Errorf(
			"生成用户 %q、session %q 的 Redis Key 失败: %w",
			userID,
			sessionID,
			err,
		)
	}
	if userMessage == nil {
		return errors.New("User Message 不能为空")
	}
	if assistantMessage == nil {
		return errors.New("Assistant Message 不能为空")
	}

	encodedUserMessage, err := json.Marshal(userMessage)
	if err != nil {
		return fmt.Errorf("序列化 User Message 失败: %w", err)
	}
	encodedAssistantMessage, err := json.Marshal(assistantMessage)
	if err != nil {
		return fmt.Errorf("序列化 Assistant Message 失败: %w", err)
	}

	// TODO 5：在一次 TxPipelined 中完成 RPUSH、LTRIM 和 EXPIRE。
	// 先 RPUSH User/Assistant 两段 JSON；再用 LTrim(ctx, key, -int64(s.maxTurns*2), -1) 只保留最近完整轮次；
	// 最后 Expire(ctx, key, s.ttl) 刷新会话过期时间。检查 TxPipelined 返回错误并用 %w 包装。
	maxMessages := int64(s.maxTurns * 2)
	_, err = s.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.RPush(ctx, key, encodedUserMessage, encodedAssistantMessage)
		pipe.LTrim(ctx, key, -maxMessages, -1)
		pipe.Expire(ctx, key, s.ttl)
		return nil
	})
	if err != nil {
		return fmt.Errorf(
			"保存用户 %q、session %q 的完整对话轮次失败: %w",
			userID,
			sessionID,
			err,
		)
	}
	return nil
}

// runTurn 重组受限历史，调用模型，并仅在成功后保存完整一轮消息。
func runTurn(
	ctx context.Context,
	agenticModel einomodel.AgenticModel,
	store historyStore,
	userID string,
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

	// TODO 6：从 Redis 读取当前用户、当前 session 的最近历史。
	// 调用 store.LoadRecent(ctx, userID, sessionID)；不要复用上一轮函数留下的内存切片。
	// 错误需要补充轮次、user ID 和 session ID 上下文。
	recent, err := store.LoadRecent(ctx, userID, sessionID)
	if err != nil {
		return "", fmt.Errorf(
			"第 %d 轮读取用户 %q、session %q 的最近历史失败: %w",
			turn,
			userID,
			sessionID,
			err,
		)
	}

	currentUserMessage := schema.UserAgenticMessage(userPrompt)
	messages := make([]*schema.AgenticMessage, 0, len(recent)+2)
	// TODO 7：按 System -> 最近历史 -> 当前 User 重组模型输入。
	// 新建 messages，先放 schema.SystemAgenticMessage(systemPrompt)，再追加 LoadRecent 返回的 User/Assistant，
	// 最后创建当前 User Message 并追加。System Message 不写 Redis，每次调用都重新注入且只出现一次。
	messages = append(messages, schema.SystemAgenticMessage(systemPrompt))
	messages = append(messages, recent...)
	messages = append(messages, currentUserMessage)
	// TODO 8：调用 generateTurn，并提取 Assistant 文本。
	// 模型或文本解析失败时直接返回，不能调用 AppendTurn，避免 Redis 中出现没有 Assistant 的半轮对话。
	message, err := generateTurn(ctx, agenticModel, messages, turn)
	// TODO 9：成功后保存当前 User 和完整 Assistant Message。
	// 调用 store.AppendTurn；保存失败时用轮次、user ID 和 session ID 包装错误；成功时返回 TODO 8 提取的文本。
	if err != nil {
		return "", err
	}
	text, err := assistantText(message)
	if err != nil {
		return "", fmt.Errorf("第 %d 轮提取 Assistant 文本失败: %w", turn, err)
	}
	err = store.AppendTurn(ctx, userID, sessionID, currentUserMessage, message)
	if err != nil {
		return "", fmt.Errorf(
			"第 %d 轮保存用户 %q、session %q 的完整对话失败: %w",
			turn,
			userID,
			sessionID,
			err,
		)
	}
	return text, nil
}

// runExercise 演示同一用户连续对话，以及不同用户使用同名 session 时的历史隔离。
func runExercise(
	ctx context.Context,
	agenticModel einomodel.AgenticModel,
	store historyStore,
) error {
	// TODO 10：使用 demoUserID 和 demoSessionID 顺序执行两轮关联对话。
	// 第一轮使用 firstUserPrompt，第二轮使用 secondUserPrompt；分别打印 User Prompt 和 Assistant 回答。
	// 第二轮必须能看到第一轮历史，而每次 runTurn 都应重新从 Redis 加载，不能依赖进程内 messages。
	firstAnswer, err := runTurn(ctx, agenticModel, store, demoUserID, demoSessionID, firstUserPrompt, 1)
	if err != nil {
		return err
	}
	fmt.Printf("User: %s\nAssistant: %s\n", firstUserPrompt, firstAnswer)
	secondAnswer, err := runTurn(ctx, agenticModel, store, demoUserID, demoSessionID, secondUserPrompt, 2)
	if err != nil {
		return err
	}
	fmt.Printf("User: %s\nAssistant: %s\n", secondUserPrompt, secondAnswer)

	// TODO 11：使用 otherDemoUserID 和相同 demoSessionID 发起隔离验证。
	// 传入 isolatedUserQuery 并打印结果。因为 Redis Key 包含 user ID，新用户不应看到 demoUserID 的两轮历史。
	// 完成后提示可使用 Redis TTL 命令观察过期时间，并返回 nil。
	isolatedAnswer, err := runTurn(ctx, agenticModel, store, otherDemoUserID, demoSessionID, isolatedUserQuery, 1)
	if err != nil {
		return err
	}
	fmt.Printf("User: %s\nAssistant: %s\n", isolatedUserQuery, isolatedAnswer)
	fmt.Printf("可使用 Redis TTL 命令观察会话 Key 的剩余过期时间（当前 TTL：%s）。\n", sessionTTL)
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

// validateAPIKey 拒绝空值和示例占位符，避免发起无效或意外的模型请求。
func validateAPIKey(key string) error {
	trimmed := strings.TrimSpace(key)
	if trimmed == "" || trimmed == "replace-with-your-api-key" {
		return errors.New("API Key 未配置")
	}
	return nil
}
