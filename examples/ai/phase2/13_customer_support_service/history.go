package main

import (
	"context"
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
	if err := validateIdentifier("user ID", userID); err != nil {
		return nil, err
	}
	if err := validateIdentifier("session ID", sessionID); err != nil {
		return nil, err
	}

	// TODO 6：构造隔离 Redis Key，读取末尾 maxTurns*2 条并严格反序列化完整 User/Assistant 轮次。
	return nil, errExerciseIncomplete
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
	if err := validateIdentifier("user ID", userID); err != nil {
		return err
	}
	if err := validateIdentifier("session ID", sessionID); err != nil {
		return err
	}
	if userMessage == nil || assistantMessage == nil {
		return fmt.Errorf("会话消息不能为空")
	}

	// TODO 7：通过 TxPipelined 原子 RPUSH 两条消息、LTRIM 到偶数上限并刷新 TTL。
	return errExerciseIncomplete
}
