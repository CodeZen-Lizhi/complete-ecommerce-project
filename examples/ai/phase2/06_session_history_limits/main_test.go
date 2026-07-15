package main

import (
	"encoding/json"
	"testing"

	"github.com/cloudwego/eino/schema"
)

// TestConversationKey 验证用户隔离 Key 的清洗、字符白名单和碰撞防护。
func TestConversationKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		userID    string
		sessionID string
		want      string
		wantErr   bool
	}{
		{name: "合法并清理空白", userID: " user-1001 ", sessionID: " checkout_help ", want: sessionKeyPrefix + "user-1001:checkout_help"},
		{name: "空用户", userID: " ", sessionID: "session", wantErr: true},
		{name: "空会话", userID: "user", sessionID: "\t", wantErr: true},
		{name: "拒绝冒号避免分段碰撞", userID: "user:admin", sessionID: "session", wantErr: true},
		{name: "拒绝内部空白", userID: "user 1", sessionID: "session", wantErr: true},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := conversationKey(test.userID, test.sessionID)
			if test.wantErr {
				if err == nil {
					t.Fatalf("conversationKey(%q, %q) 未返回预期错误，结果为 %q", test.userID, test.sessionID, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("conversationKey(%q, %q) 返回错误: %v", test.userID, test.sessionID, err)
			}
			if got != test.want {
				t.Fatalf("conversationKey(%q, %q) = %q，期望 %q", test.userID, test.sessionID, got, test.want)
			}
		})
	}
}

// TestDecodeHistoryValues 验证 Redis 历史只接受完整、非空且可解码的消息对。
func TestDecodeHistoryValues(t *testing.T) {
	t.Parallel()

	userJSON := mustMarshalMessage(t, schema.UserAgenticMessage("hello"))
	assistantJSON := mustMarshalMessage(t, schema.UserAgenticMessage("world"))

	tests := []struct {
		name    string
		values  []string
		wantLen int
		wantErr bool
	}{
		{name: "空历史", values: nil, wantLen: 0},
		{name: "完整一轮", values: []string{userJSON, assistantJSON}, wantLen: 2},
		{name: "奇数条历史", values: []string{userJSON}, wantErr: true},
		{name: "空元素", values: []string{"", assistantJSON}, wantErr: true},
		{name: "损坏 JSON", values: []string{"{", assistantJSON}, wantErr: true},
		{name: "null 消息", values: []string{"null", assistantJSON}, wantErr: true},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			messages, err := decodeHistoryValues("user-1001", "session-1", test.values)
			if test.wantErr {
				if err == nil {
					t.Fatalf("decodeHistoryValues(%v) 未返回预期错误", test.values)
				}
				return
			}
			if err != nil {
				t.Fatalf("decodeHistoryValues(%v) 返回错误: %v", test.values, err)
			}
			if len(messages) != test.wantLen {
				t.Fatalf("消息数量为 %d，期望 %d", len(messages), test.wantLen)
			}
		})
	}
}

// mustMarshalMessage 序列化测试消息，失败时立即终止当前测试。
func mustMarshalMessage(t *testing.T, message *schema.AgenticMessage) string {
	t.Helper()

	encoded, err := json.Marshal(message)
	if err != nil {
		t.Fatalf("序列化测试消息失败: %v", err)
	}
	return string(encoded)
}
