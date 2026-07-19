package main

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	einomodel "github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/redis/go-redis/v9"
)

type fakeHistoryStore struct {
	history           []*schema.AgenticMessage
	loadErr           error
	appendErr         error
	appendCalls       int
	appendedUser      *schema.AgenticMessage
	appendedAssistant *schema.AgenticMessage
}

// LoadRecent 返回测试预置的会话历史或错误。
func (s *fakeHistoryStore) LoadRecent(context.Context, string, string) ([]*schema.AgenticMessage, error) {
	return s.history, s.loadErr
}

// AppendTurn 记录待保存的完整轮次，并返回测试预置错误。
func (s *fakeHistoryStore) AppendTurn(
	_ context.Context,
	_ string,
	_ string,
	userMessage *schema.AgenticMessage,
	assistantMessage *schema.AgenticMessage,
) error {
	s.appendCalls++
	s.appendedUser = userMessage
	s.appendedAssistant = assistantMessage
	return s.appendErr
}

type fakeAgenticModel struct {
	response *schema.AgenticMessage
	err      error
	input    []*schema.AgenticMessage
}

// Generate 记录模型输入并返回测试预置响应。
func (m *fakeAgenticModel) Generate(
	_ context.Context,
	input []*schema.AgenticMessage,
	_ ...einomodel.Option,
) (*schema.AgenticMessage, error) {
	m.input = input
	return m.response, m.err
}

// Stream 满足 AgenticModel 接口；本组测试不会调用流式方法。
func (m *fakeAgenticModel) Stream(
	context.Context,
	[]*schema.AgenticMessage,
	...einomodel.Option,
) (*schema.StreamReader[*schema.AgenticMessage], error) {
	return nil, errors.New("测试未实现 Stream")
}

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
	assistantJSON := mustMarshalMessage(t, newAssistantMessage("world"))
	secondUserJSON := mustMarshalMessage(t, schema.UserAgenticMessage("second user"))
	firstAssistantJSON := mustMarshalMessage(t, newAssistantMessage("first assistant"))

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
		{name: "User 后仍是 User", values: []string{userJSON, secondUserJSON}, wantErr: true},
		{name: "第一条不是 User", values: []string{firstAssistantJSON, assistantJSON}, wantErr: true},
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

// TestRunTurn 验证历史重组、文本解析和完整轮次保存顺序。
func TestRunTurn(t *testing.T) {
	t.Parallel()

	history := []*schema.AgenticMessage{
		schema.UserAgenticMessage("历史问题"),
		newAssistantMessage("历史回答"),
	}
	response := newAssistantMessage("当前回答")
	store := &fakeHistoryStore{history: history}
	model := &fakeAgenticModel{response: response}

	text, err := runTurn(context.Background(), model, store, "user-1", "session-1", "当前问题", 1)
	if err != nil {
		t.Fatalf("runTurn 返回错误: %v", err)
	}
	if text != "当前回答" {
		t.Fatalf("Assistant 文本 = %q，期望 %q", text, "当前回答")
	}
	if len(model.input) != len(history)+2 {
		t.Fatalf("模型输入消息数 = %d，期望 %d", len(model.input), len(history)+2)
	}
	if model.input[0].Role != schema.AgenticRoleTypeSystem || model.input[len(model.input)-1].Role != schema.AgenticRoleTypeUser {
		t.Fatalf("模型输入顺序不是 System -> History -> User")
	}
	if store.appendCalls != 1 || store.appendedAssistant != response {
		t.Fatalf("成功响应应保存且只保存一次完整轮次")
	}
}

// TestRunTurnRejectsInvalidAssistantBeforeSaving 验证无有效 Assistant 文本时不会污染历史。
func TestRunTurnRejectsInvalidAssistantBeforeSaving(t *testing.T) {
	t.Parallel()

	store := &fakeHistoryStore{}
	model := &fakeAgenticModel{response: schema.UserAgenticMessage("错误角色内容")}

	_, err := runTurn(context.Background(), model, store, "user-1", "session-1", "问题", 1)
	if err == nil || !strings.Contains(err.Error(), "提取 Assistant 文本失败") {
		t.Fatalf("无有效 Assistant 文本应返回明确错误，实际为: %v", err)
	}
	if store.appendCalls != 0 {
		t.Fatalf("解析失败时不应保存历史，实际保存 %d 次", store.appendCalls)
	}
}

// TestRunTurnDoesNotSaveWhenModelFails 验证模型错误时不会写入任何会话历史。
func TestRunTurnDoesNotSaveWhenModelFails(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("模型调用失败")
	store := &fakeHistoryStore{}
	model := &fakeAgenticModel{err: wantErr}

	_, err := runTurn(context.Background(), model, store, "user-1", "session-1", "问题", 1)
	if !errors.Is(err, wantErr) {
		t.Fatalf("模型错误链未保留: %v", err)
	}
	if store.appendCalls != 0 {
		t.Fatalf("模型失败时不应保存历史，实际保存 %d 次", store.appendCalls)
	}
}

// TestRunTurnPreservesAppendErrorContext 验证保存失败时保留底层错误和轮次、用户、会话上下文。
func TestRunTurnPreservesAppendErrorContext(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("Redis 写入失败")
	store := &fakeHistoryStore{appendErr: wantErr}
	model := &fakeAgenticModel{response: newAssistantMessage("回答")}

	_, err := runTurn(context.Background(), model, store, "user-1", "session-1", "问题", 3)
	if !errors.Is(err, wantErr) {
		t.Fatalf("保存错误链未保留: %v", err)
	}
	for _, expected := range []string{"第 3 轮", `user-1`, `session-1`} {
		if !strings.Contains(err.Error(), expected) {
			t.Fatalf("保存错误缺少上下文 %q: %v", expected, err)
		}
	}
	if store.appendCalls != 1 {
		t.Fatalf("保存失败时 AppendTurn 调用次数 = %d，期望 1", store.appendCalls)
	}
}

// TestAppendTurnRejectsInvalidRoles 验证写 Redis 前必须先满足 User/Assistant 角色契约。
func TestAppendTurnRejectsInvalidRoles(t *testing.T) {
	t.Parallel()

	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:0"})
	t.Cleanup(func() {
		if err := client.Close(); err != nil {
			t.Errorf("关闭测试 Redis Client 失败: %v", err)
		}
	})
	store, err := newRedisHistoryStore(client, 1, time.Minute)
	if err != nil {
		t.Fatalf("创建测试历史存储失败: %v", err)
	}

	tests := []struct {
		name      string
		user      *schema.AgenticMessage
		assistant *schema.AgenticMessage
		wantText  string
	}{
		{name: "User 角色错误", user: newAssistantMessage("错误"), assistant: newAssistantMessage("回答"), wantText: "User Message 角色"},
		{name: "Assistant 角色错误", user: schema.UserAgenticMessage("问题"), assistant: schema.UserAgenticMessage("错误"), wantText: "Assistant Message 角色"},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			err := store.AppendTurn(context.Background(), "user-1", "session-1", test.user, test.assistant)
			if err == nil || !strings.Contains(err.Error(), test.wantText) {
				t.Fatalf("期望角色错误包含 %q，实际为: %v", test.wantText, err)
			}
		})
	}
}

// newAssistantMessage 创建只包含可展示文本块的 Assistant 消息。
func newAssistantMessage(text string) *schema.AgenticMessage {
	return &schema.AgenticMessage{
		Role: schema.AgenticRoleTypeAssistant,
		ContentBlocks: []*schema.ContentBlock{
			schema.NewContentBlock(&schema.AssistantGenText{Text: text}),
		},
	}
}
