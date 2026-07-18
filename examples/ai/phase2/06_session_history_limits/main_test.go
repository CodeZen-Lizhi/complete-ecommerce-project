package main

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	einomodel "github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
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

// newAssistantMessage 创建只包含可展示文本块的 Assistant 消息。
func newAssistantMessage(text string) *schema.AgenticMessage {
	return &schema.AgenticMessage{
		Role: schema.AgenticRoleTypeAssistant,
		ContentBlocks: []*schema.ContentBlock{
			schema.NewContentBlock(&schema.AssistantGenText{Text: text}),
		},
	}
}
