package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/redis/go-redis/v9"
)

const historyKeyPrefix = "ai:phase2:customer-support:"

// historyStore 隔离读写指定用户会话的完整 User/Assistant 历史轮次。
type historyStore interface {
	// Load 读取指定用户和会话的完整历史消息对；空历史不视为失败。
	Load(context.Context, string, string) ([]*schema.Message, error)
	// AppendTurn 原子写入一条完整 User/Assistant 会话轮次。
	AppendTurn(context.Context, string, string, *schema.Message, *schema.Message) error
}

// redisHistoryStore 用 Redis 保存有 TTL 和轮次上限的会话历史。
type redisHistoryStore struct {
	client   redis.Cmdable
	maxTurns int
	ttl      time.Duration
}

// newRedisHistoryStore 校验 Redis 依赖和会话边界后创建唯一历史 Store。
func newRedisHistoryStore(client redis.Cmdable, maxTurns int, ttl time.Duration) (historyStore, error) {
	if client == nil {
		return nil, fmt.Errorf("Redis Client 不能为空")
	}
	if maxTurns <= 0 || ttl <= 0 {
		return nil, fmt.Errorf("会话历史上限或 TTL 无效")
	}
	return &redisHistoryStore{client: client, maxTurns: maxTurns, ttl: ttl}, nil
}

// Load 读取当前用户当前会话最近的完整历史消息；空历史返回空切片而非错误。
func (store *redisHistoryStore) Load(ctx context.Context, userID string, sessionID string) ([]*schema.Message, error) {
	if ctx == nil {
		return nil, fmt.Errorf("历史读取 Context 不能为空")
	}
	if store == nil || store.client == nil {
		return nil, fmt.Errorf("Redis 历史 Store 未初始化")
	}
	if err := validateIdentifier("user ID", userID); err != nil {
		return nil, err
	}
	if err := validateIdentifier("session ID", sessionID); err != nil {
		return nil, err
	}
	redisKey := historyKeyPrefix + userID + ":" + sessionID
	maxMessages := int64(store.maxTurns) * 2
	values, err := store.client.LRange(ctx, redisKey, -maxMessages, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("读取 Redis 会话历史失败: %w", err)
	}
	if len(values) == 0 {
		return []*schema.Message{}, nil
	}
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("Redis 会话历史包含 %d 条消息，不是完整的 User/Assistant 轮次", len(values))
	}
	messages := make([]*schema.Message, 0, len(values))
	for index, value := range values {
		var message *schema.Message
		if err := json.Unmarshal([]byte(value), &message); err != nil {
			return nil, fmt.Errorf("反序列化第 %d 条历史消息失败: %w", index, err)
		}
		if message == nil {
			return nil, fmt.Errorf("第 %d 条历史消息为 null", index)
		}

		expectedRole := schema.User
		if index%2 == 1 {
			expectedRole = schema.Assistant
		}
		if message.Role != expectedRole {
			return nil, fmt.Errorf(
				"第 %d 条消息角色为 %q，期望 %q",
				index,
				message.Role,
				expectedRole,
			)
		}
		messages = append(messages, message)
	}

	// TODO 6：构造隔离 Redis Key，读取末尾 maxTurns*2 条并严格反序列化完整 User/Assistant 轮次。
	return messages, nil
}

// AppendTurn 原子写入一整轮 User/Assistant、裁剪历史并刷新 TTL，任何失败都不留下半轮数据。
func (store *redisHistoryStore) AppendTurn(
	ctx context.Context,
	userID string,
	sessionID string,
	userMessage *schema.Message,
	assistantMessage *schema.Message,
) error {
	if ctx == nil {
		return fmt.Errorf("历史写入 Context 不能为空")
	}
	if store == nil || store.client == nil {
		return fmt.Errorf("Redis 历史 Store 未初始化")
	}
	if err := validateIdentifier("user ID", userID); err != nil {
		return err
	}
	if err := validateIdentifier("session ID", sessionID); err != nil {
		return err
	}
	if userMessage == nil || assistantMessage == nil {
		return fmt.Errorf("会话消息不能为空")
	}
	if userMessage.Role != schema.User {
		return fmt.Errorf("User Message 角色为 %q，期望 %q", userMessage.Role, schema.User)
	}
	if assistantMessage.Role != schema.Assistant {
		return fmt.Errorf("Assistant Message 角色为 %q，期望 %q", assistantMessage.Role, schema.Assistant)
	}

	encodedUserMessage, err := json.Marshal(userMessage)
	if err != nil {
		return fmt.Errorf("序列化 User Message 失败: %w", err)
	}
	encodedAssistantMessage, err := json.Marshal(assistantMessage)
	if err != nil {
		return fmt.Errorf("序列化 Assistant Message 失败: %w", err)
	}

	redisKey := historyKeyPrefix + userID + ":" + sessionID
	maxMessages := int64(store.maxTurns) * 2
	_, err = store.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.RPush(ctx, redisKey, encodedUserMessage, encodedAssistantMessage)
		pipe.LTrim(ctx, redisKey, -maxMessages, -1)
		pipe.Expire(ctx, redisKey, store.ttl)
		return nil
	})
	if err != nil {
		return fmt.Errorf("保存完整 Redis 会话轮次失败: %w", err)
	}
	// TODO 7：通过 TxPipelined 原子 RPUSH 两条消息、LTRIM 到偶数上限并刷新 TTL。
	return nil
}
