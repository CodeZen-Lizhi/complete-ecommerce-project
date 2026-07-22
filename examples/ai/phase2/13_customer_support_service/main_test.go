package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

// TestLoadExerciseConfig 验证环境配置在创建外部客户端前完成边界校验。
func TestLoadExerciseConfig(t *testing.T) {
	values := map[string]string{
		"OPENAI_API_KEY": "local-test-key",
	}
	config, err := loadExerciseConfig(func(key string) string { return values[key] })
	if err != nil {
		t.Fatalf("加载默认练习配置失败: %v", err)
	}
	if config.ListenAddr != defaultListenAddr || config.HistoryTurns != defaultHistoryTurns {
		t.Fatalf("默认配置不符合预期: %+v", config)
	}

	for _, invalid := range []struct {
		name  string
		key   string
		value string
	}{
		{name: "non loopback address", key: "AI_DEMO_LISTEN_ADDR", value: "0.0.0.0:8093"},
		{name: "nan rate", key: "AI_DEMO_CALLS_PER_SECOND", value: "NaN"},
		{name: "infinite rate", key: "AI_DEMO_CALLS_PER_SECOND", value: "+Inf"},
	} {
		t.Run(invalid.name, func(t *testing.T) {
			values[invalid.key] = invalid.value
			t.Cleanup(func() { delete(values, invalid.key) })
			if _, err := loadExerciseConfig(func(key string) string { return values[key] }); err == nil {
				t.Fatalf("%s=%q 应被拒绝", invalid.key, invalid.value)
			}
		})
	}
}

// TestDefaultBusinessKnowledge 验证内置知识包含且只包含四个受支持意图。
func TestDefaultBusinessKnowledge(t *testing.T) {
	knowledge := defaultBusinessKnowledge()
	if err := validateBusinessKnowledge(knowledge); err != nil {
		t.Fatalf("校验内置业务知识失败: %v", err)
	}
	if _, err := knowledge.contextFor(intentDeliveryReturn); err != nil {
		t.Fatalf("读取配送退换知识失败: %v", err)
	}
	if _, err := knowledge.contextFor(customerIntent("unknown")); err == nil {
		t.Fatal("未知意图应被拒绝")
	}
}

// TestBuildClassificationResponseFormat 验证分类阶段使用严格、封闭且字段完整的 JSON Schema。
func TestBuildClassificationResponseFormat(t *testing.T) {
	format, err := buildClassificationResponseFormat()
	if err != nil {
		t.Fatalf("构造分类 JSON Schema 失败: %v", err)
	}
	if format == nil || format.Type != openai.ChatCompletionResponseFormatTypeJSONSchema {
		t.Fatal("分类响应格式必须使用原生 JSON Schema")
	}
	if format.JSONSchema == nil || !format.JSONSchema.Strict || format.JSONSchema.JSONSchema == nil {
		t.Fatal("分类响应格式必须启用 Strict 并提供 JSON Schema")
	}

	rawSchema, err := json.Marshal(format.JSONSchema.JSONSchema)
	if err != nil {
		t.Fatalf("序列化分类 JSON Schema 失败: %v", err)
	}
	schemaText := string(rawSchema)
	for _, expected := range []string{
		`"additionalProperties":false`,
		`"required":["intent","response_style","requires_handoff"]`,
		`"product_advice"`,
		`"delivery_return"`,
		`"after_sales"`,
		`"general"`,
		`"concise"`,
		`"guided"`,
		`"empathetic"`,
	} {
		if !strings.Contains(schemaText, expected) {
			t.Fatalf("分类 JSON Schema 缺少约束 %s: %s", expected, schemaText)
		}
	}
}

// TestClassificationSystemMessage 验证分类 System Message 使用 System 角色并覆盖固定分类契约。
func TestClassificationSystemMessage(t *testing.T) {
	message, err := classificationSystemMessage()
	if err != nil {
		t.Fatalf("构造分类 System Message 失败: %v", err)
	}
	if message == nil || message.Role != schema.System {
		t.Fatal("分类消息必须使用 System 角色")
	}
	if strings.TrimSpace(message.Content) == "" {
		t.Fatal("分类 System Message 不能为空")
	}
	for _, expected := range []string{
		"product_advice",
		"delivery_return",
		"after_sales",
		"general",
		"concise",
		"guided",
		"empathetic",
		"requires_handoff",
		"JSON Schema",
	} {
		if !strings.Contains(message.Content, expected) {
			t.Fatalf("分类 System Message 缺少约束 %q: %s", expected, message.Content)
		}
	}
}

// TestClassificationFewShotMessages 验证三个分类样例按 User/Assistant 顺序提供合法分类 JSON。
func TestClassificationFewShotMessages(t *testing.T) {
	messages, err := classificationFewShotMessages()
	if err != nil {
		t.Fatalf("构造分类 Few-shot 消息失败: %v", err)
	}
	expectedResults := []classification{
		{Intent: intentProductAdvice, ResponseStyle: styleGuided, RequiresHandoff: false},
		{Intent: intentDeliveryReturn, ResponseStyle: styleConcise, RequiresHandoff: false},
		{Intent: intentAfterSales, ResponseStyle: styleEmpathetic, RequiresHandoff: true},
	}
	if len(messages) != len(expectedResults)*2 {
		t.Fatalf("Few-shot 消息数量 = %d，期望 %d", len(messages), len(expectedResults)*2)
	}

	for index, expected := range expectedResults {
		userMessage := messages[index*2]
		assistantMessage := messages[index*2+1]
		if userMessage == nil || userMessage.Role != schema.User || strings.TrimSpace(userMessage.Content) == "" {
			t.Fatalf("第 %d 个 Few-shot User 消息无效: %+v", index, userMessage)
		}
		if assistantMessage == nil || assistantMessage.Role != schema.Assistant {
			t.Fatalf("第 %d 个 Few-shot Assistant 消息无效: %+v", index, assistantMessage)
		}

		var actual classification
		if err := json.Unmarshal([]byte(assistantMessage.Content), &actual); err != nil {
			t.Fatalf("第 %d 个 Few-shot Assistant JSON 无效: %v", index, err)
		}
		if actual != expected {
			t.Fatalf("第 %d 个 Few-shot 分类 = %+v，期望 %+v", index, actual, expected)
		}
	}
}

// TestBuildAnswerMessages 验证回答模板按 System、历史、当前 User 的顺序展开 typed 消息历史。
func TestBuildAnswerMessages(t *testing.T) {
	history := []*schema.Message{
		schema.UserMessage("之前的问题"),
		schema.AssistantMessage("之前的回答", nil),
	}
	classificationResult := classification{
		Intent:          intentDeliveryReturn,
		ResponseStyle:   styleGuided,
		RequiresHandoff: false,
	}

	messages, err := buildAnswerMessages(
		context.Background(),
		classificationResult,
		"配送状态以订单页为准。",
		history,
		"我的订单什么时候到？",
	)
	if err != nil {
		t.Fatalf("构造回答消息失败: %v", err)
	}
	if len(messages) != len(history)+2 {
		t.Fatalf("回答消息数量 = %d，期望 %d", len(messages), len(history)+2)
	}
	if messages[0] == nil || messages[0].Role != schema.System {
		t.Fatalf("第一条消息必须是 System: %+v", messages[0])
	}
	for _, expected := range []string{
		string(intentDeliveryReturn),
		string(styleGuided),
		"配送状态以订单页为准。",
	} {
		if !strings.Contains(messages[0].Content, expected) {
			t.Fatalf("回答 System Message 缺少 %q: %s", expected, messages[0].Content)
		}
	}
	for index, expected := range history {
		actual := messages[index+1]
		if actual == nil || actual.Role != expected.Role || actual.Content != expected.Content {
			t.Fatalf("第 %d 条历史消息未按原顺序展开: %+v", index, actual)
		}
	}
	if messages[3] == nil || messages[3].Role != schema.User || messages[3].Content != "我的订单什么时候到？" {
		t.Fatalf("最后一条消息必须是当前 User: %+v", messages[3])
	}

	emptyHistoryMessages, err := buildAnswerMessages(
		context.Background(),
		classificationResult,
		"配送状态以订单页为准。",
		nil,
		"怎么查询物流？",
	)
	if err != nil {
		t.Fatalf("空历史应允许格式化: %v", err)
	}
	if len(emptyHistoryMessages) != 2 || emptyHistoryMessages[0].Role != schema.System || emptyHistoryMessages[1].Role != schema.User {
		t.Fatalf("空历史消息顺序不正确: %+v", emptyHistoryMessages)
	}
}

// TestRedisHistoryStore 验证 Redis 历史按完整轮次读取，并原子追加、裁剪和刷新 TTL。
func TestRedisHistoryStore(t *testing.T) {
	userMessage := schema.UserMessage("之前的问题")
	assistantMessage := schema.AssistantMessage("之前的回答", nil)
	encodedUserMessage, err := json.Marshal(userMessage)
	if err != nil {
		t.Fatalf("编码测试 User Message 失败: %v", err)
	}
	encodedAssistantMessage, err := json.Marshal(assistantMessage)
	if err != nil {
		t.Fatalf("编码测试 Assistant Message 失败: %v", err)
	}

	hook := &redisHistoryHook{
		lrangeValues: []string{string(encodedUserMessage), string(encodedAssistantMessage)},
	}
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:0"})
	client.AddHook(hook)
	t.Cleanup(func() { _ = client.Close() })

	store, err := newRedisHistoryStore(client, 2, time.Minute)
	if err != nil {
		t.Fatalf("创建 Redis 历史 Store 失败: %v", err)
	}
	messages, err := store.Load(context.Background(), "learner-1001", "checkout-help")
	if err != nil {
		t.Fatalf("读取 Redis 历史失败: %v", err)
	}
	if len(messages) != 2 || messages[0].Role != schema.User || messages[1].Role != schema.Assistant {
		t.Fatalf("Redis 历史轮次不完整: %+v", messages)
	}
	if len(hook.processArgs) != 4 || hook.processArgs[0] != "lrange" || hook.processArgs[1] != historyKeyPrefix+"learner-1001:checkout-help" {
		t.Fatalf("LRange 参数不正确: %+v", hook.processArgs)
	}
	if start, ok := hook.processArgs[2].(int64); !ok || start != -4 {
		t.Fatalf("LRange 起始下标 = %#v，期望 -4", hook.processArgs[2])
	}
	if stop, ok := hook.processArgs[3].(int64); !ok || stop != -1 {
		t.Fatalf("LRange 结束下标 = %#v，期望 -1", hook.processArgs[3])
	}

	if err := store.AppendTurn(context.Background(), "learner-1001", "checkout-help", userMessage, assistantMessage); err != nil {
		t.Fatalf("追加 Redis 历史失败: %v", err)
	}
	expectedCommands := []string{"multi", "rpush", "ltrim", "expire", "exec"}
	if len(hook.pipelineArgs) != len(expectedCommands) {
		t.Fatalf("事务命令数量 = %d，期望 %d: %+v", len(hook.pipelineArgs), len(expectedCommands), hook.pipelineArgs)
	}
	for index, expectedCommand := range expectedCommands {
		if len(hook.pipelineArgs[index]) == 0 || hook.pipelineArgs[index][0] != expectedCommand {
			t.Fatalf("第 %d 条事务命令不正确: %+v", index, hook.pipelineArgs[index])
		}
	}
	for index, value := range hook.pipelineArgs[1][2:] {
		raw, ok := value.([]byte)
		if !ok {
			t.Fatalf("RPUSH 第 %d 条消息不是 JSON 字节: %T", index, value)
		}
		var message schema.Message
		if err := json.Unmarshal(raw, &message); err != nil {
			t.Fatalf("RPUSH 第 %d 条消息不是合法 JSON: %v", index, err)
		}
	}

	hook.lrangeValues = []string{string(encodedUserMessage)}
	if _, err := store.Load(context.Background(), "learner-1001", "checkout-help"); err == nil {
		t.Fatal("奇数条 Redis 历史必须被拒绝")
	}
}

// TestIncompleteExerciseStopsBeforeProvider 验证未完成核心步骤不触发 Provider 或历史写入。
func TestIncompleteExerciseStopsBeforeProvider(t *testing.T) {
	service, provider, history := newTestCustomerSupportService(t)
	err := service.Stream(
		context.Background(),
		"learner-1001",
		chatRequest{SessionID: "checkout-help", Message: "这款耳机多久能到？"},
		&recordingEmitter{},
	)
	if !errors.Is(err, errExerciseIncomplete) {
		t.Fatalf("未完成服务错误 = %v，期望 errExerciseIncomplete", err)
	}
	if provider.generateCalls != 0 || provider.streamCalls != 0 {
		t.Fatalf("未完成练习不应调用 Provider，generate=%d stream=%d", provider.generateCalls, provider.streamCalls)
	}
	if history.appendCalls != 0 {
		t.Fatalf("未完成练习不应写入历史，append=%d", history.appendCalls)
	}
}

// TestHTTPRejectsForgedBodyIdentity 验证请求体不能伪造用户身份，且未完成路径返回安全 JSON 错误。
func TestHTTPRejectsForgedBodyIdentity(t *testing.T) {
	service, provider, _ := newTestCustomerSupportService(t)
	router := newCustomerSupportRouter(service)

	for _, testCase := range []struct {
		name        string
		body        string
		withHeader  bool
		contentType string
		wantStatus  int
		wantCode    string
	}{
		{
			name:        "missing identity",
			body:        `{"session_id":"checkout-help","message":"你好"}`,
			contentType: "application/json",
			wantStatus:  http.StatusUnauthorized,
			wantCode:    "unauthenticated",
		},
		{
			name:        "forged body identity",
			body:        `{"session_id":"checkout-help","message":"你好","user_id":"admin"}`,
			withHeader:  true,
			contentType: "application/json",
			wantStatus:  http.StatusBadRequest,
			wantCode:    "invalid_request",
		},
		{
			name:        "incomplete exercise",
			body:        `{"session_id":"checkout-help","message":"你好"}`,
			withHeader:  true,
			contentType: "application/json",
			wantStatus:  http.StatusServiceUnavailable,
			wantCode:    "exercise_incomplete",
		},
		{
			name:        "unsupported media type",
			body:        `{"session_id":"checkout-help","message":"你好"}`,
			withHeader:  true,
			contentType: "text/plain",
			wantStatus:  http.StatusUnsupportedMediaType,
			wantCode:    "unsupported_media_type",
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/api/ai/chat/stream", strings.NewReader(testCase.body))
			request.Header.Set("Content-Type", testCase.contentType)
			if testCase.withHeader {
				request.Header.Set(demoUserIDHeader, "learner-1001")
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			if response.Code != testCase.wantStatus {
				t.Fatalf("状态码 = %d，期望 %d，响应=%s", response.Code, testCase.wantStatus, response.Body.String())
			}
			var payload apiError
			if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
				t.Fatalf("解码 JSON 错误响应失败: %v", err)
			}
			if payload.Code != testCase.wantCode {
				t.Fatalf("错误码 = %q，期望 %q", payload.Code, testCase.wantCode)
			}
		})
	}
	if provider.generateCalls != 0 || provider.streamCalls != 0 {
		t.Fatalf("HTTP 校验和未完成路径不应调用 Provider，generate=%d stream=%d", provider.generateCalls, provider.streamCalls)
	}
}

// TestCustomerSupportServiceHappyPathAfterCompletion 在 TODO 完成后验证完整流、历史隔离和事件顺序。
func TestCustomerSupportServiceHappyPathAfterCompletion(t *testing.T) {
	service, provider, history := newTestCustomerSupportService(t)
	stream := &fakeMessageStream{steps: []streamStep{
		{chunk: modelChunk{Content: "您好，"}},
		{chunk: modelChunk{Content: "请在订单页查看配送进度。"}},
		{chunk: modelChunk{Usage: tokenUsage{PromptTokens: 12, CompletionTokens: 8, TotalTokens: 20, Available: true}}},
	}}
	provider.generateResults = []modelResult{{Content: validClassificationJSON}}
	provider.streamResults = []messageStream{stream}
	emitter := &recordingEmitter{
		beforeSend: func(event string) {
			if event == "done" && history.appendCalls != 1 {
				t.Errorf("发送 done 前应已提交一轮历史，append=%d", history.appendCalls)
			}
		},
	}

	err := service.Stream(
		context.Background(),
		"learner-1001",
		chatRequest{SessionID: "checkout-help", Message: "这款耳机多久能到？"},
		emitter,
	)
	requireCompletedExercise(t, err)
	if err != nil {
		t.Fatalf("完整客服流程失败: %v", err)
	}
	if provider.generateCalls != 1 || provider.streamCalls != 1 {
		t.Fatalf("Provider 调用次数不正确，generate=%d stream=%d", provider.generateCalls, provider.streamCalls)
	}
	if history.loadCalls != 1 || history.appendCalls != 1 || history.lastLoadUserID != "learner-1001" || history.lastAppendSessionID != "checkout-help" {
		t.Fatalf("历史隔离或提交参数不正确: %+v", history)
	}
	if history.lastUserMessage == nil || history.lastUserMessage.Content != "这款耳机多久能到？" || history.lastAssistantMessage == nil || history.lastAssistantMessage.Content != "您好，请在订单页查看配送进度。" {
		t.Fatalf("提交的完整会话轮次不正确: user=%+v assistant=%+v", history.lastUserMessage, history.lastAssistantMessage)
	}
	if got := strings.Join(emitter.eventNames(), ","); got != "meta,delta,delta,done" {
		t.Fatalf("SSE 事件顺序 = %q，期望 meta,delta,delta,done", got)
	}
	meta, ok := emitter.events[0].payload.(metaEvent)
	if !ok || meta.Intent != intentDeliveryReturn || meta.ResponseStyle != styleGuided || meta.RequiresHandoff {
		t.Fatalf("meta 事件不正确: %#v", emitter.events[0].payload)
	}
	if !stream.closed {
		t.Fatal("完整流程结束后必须关闭回答流")
	}
}

// TestCustomerSupportServiceDoesNotCommitOnStreamErrorAfterCompletion 验证首个 delta 后的流错误不会重试或写历史。
func TestCustomerSupportServiceDoesNotCommitOnStreamErrorAfterCompletion(t *testing.T) {
	service, provider, history := newTestCustomerSupportService(t)
	stream := &fakeMessageStream{steps: []streamStep{
		{chunk: modelChunk{Content: "部分回答"}},
		{err: errors.New("模拟流接收失败")},
	}}
	provider.generateResults = []modelResult{{Content: validClassificationJSON}}
	provider.streamResults = []messageStream{stream}
	emitter := &recordingEmitter{}

	err := service.Stream(
		context.Background(),
		"learner-1001",
		chatRequest{SessionID: "checkout-help", Message: "配送有延迟怎么办？"},
		emitter,
	)
	requireCompletedExercise(t, err)
	if err == nil {
		t.Fatal("流中途错误必须返回失败")
	}
	if provider.streamCalls != 1 {
		t.Fatalf("首个 delta 后不得重建流，stream=%d", provider.streamCalls)
	}
	if history.appendCalls != 0 {
		t.Fatalf("流中途错误不得写入历史，append=%d", history.appendCalls)
	}
	if got := strings.Join(emitter.eventNames(), ","); got != "meta,delta" {
		t.Fatalf("流错误前事件顺序 = %q，期望 meta,delta", got)
	}
	if !stream.closed {
		t.Fatal("流中途错误后必须关闭回答流")
	}
}

// TestCustomerSupportServiceDoesNotCommitOnHistoryFailureAfterCompletion 验证 Redis 提交失败不会发送 done。
func TestCustomerSupportServiceDoesNotCommitOnHistoryFailureAfterCompletion(t *testing.T) {
	service, provider, history := newTestCustomerSupportService(t)
	provider.generateResults = []modelResult{{Content: validClassificationJSON}}
	provider.streamResults = []messageStream{&fakeMessageStream{steps: []streamStep{{chunk: modelChunk{Content: "完整回答"}}}}}
	history.appendErr = errors.New("模拟 Redis 提交失败")
	emitter := &recordingEmitter{}

	err := service.Stream(
		context.Background(),
		"learner-1001",
		chatRequest{SessionID: "checkout-help", Message: "售后如何处理？"},
		emitter,
	)
	requireCompletedExercise(t, err)
	if err == nil {
		t.Fatal("历史提交失败必须返回错误")
	}
	if history.appendCalls != 1 {
		t.Fatalf("历史提交失败仍应只尝试一次，append=%d", history.appendCalls)
	}
	if got := strings.Join(emitter.eventNames(), ","); got != "meta,delta" {
		t.Fatalf("历史失败不得发送 done，事件=%q", got)
	}
}

// TestClassificationValidationAfterCompletion 验证严格分类解析拒绝未知字段和不完整结果。
func TestClassificationValidationAfterCompletion(t *testing.T) {
	_, err := decodeAndValidateClassification([]byte(`{"intent":"delivery_return","response_style":"guided","requires_handoff":false,"extra":true}`))
	requireCompletedExercise(t, err)
	if err == nil {
		t.Fatal("未知 JSON 字段必须被拒绝")
	}

	_, err = decodeAndValidateClassification([]byte(`{"intent":"delivery_return","response_style":"guided"}`))
	if err == nil {
		t.Fatal("缺失 requires_handoff 必须被拒绝")
	}
}

// TestGovernedGenerateRetriesTemporaryFailureAfterCompletion 验证临时网络错误按上限重试并记录尝试次数。
func TestGovernedGenerateRetriesTemporaryFailureAfterCompletion(t *testing.T) {
	provider := &fakeChatProvider{generateErrors: []error{temporaryNetworkError{}, temporaryNetworkError{}}}
	_, metrics, err := governedGenerate(
		context.Background(),
		provider,
		governanceConfig{
			Timeout:        time.Second,
			MaxRetries:     1,
			InitialBackoff: time.Nanosecond,
			Limiter:        rate.NewLimiter(rate.Inf, 1),
		},
		modelCall{Messages: []*schema.Message{schema.UserMessage("测试重试")}},
	)
	requireCompletedExercise(t, err)
	if err == nil {
		t.Fatal("重试耗尽必须返回最后错误")
	}
	if provider.generateCalls != 2 || metrics.Attempts != 2 {
		t.Fatalf("重试次数不正确，provider=%d metrics=%d", provider.generateCalls, metrics.Attempts)
	}
}

// TestHTTPEventOrderAfterCompletion 验证完成实现后 HTTP 层输出 meta、delta、done 的 SSE 顺序。
func TestHTTPEventOrderAfterCompletion(t *testing.T) {
	service, provider, history := newTestCustomerSupportService(t)
	provider.generateResults = []modelResult{{Content: validClassificationJSON}}
	provider.streamResults = []messageStream{&fakeMessageStream{steps: []streamStep{{chunk: modelChunk{Content: "请查看订单页。"}}}}}
	router := newCustomerSupportRouter(service)
	request := httptest.NewRequest(http.MethodPost, "/api/ai/chat/stream", strings.NewReader(`{"session_id":"checkout-help","message":"配送进度"}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set(demoUserIDHeader, "learner-1001")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code == http.StatusServiceUnavailable {
		var payload apiError
		if err := json.Unmarshal(response.Body.Bytes(), &payload); err == nil && payload.Code == "exercise_incomplete" {
			t.Skip("核心 TODO 尚未完成；完成后自动启用 HTTP SSE 顺序断言")
		}
	}
	if response.Code != http.StatusOK {
		t.Fatalf("完成实现后的 HTTP 状态 = %d，响应=%s", response.Code, response.Body.String())
	}
	if got := sseEventNames(response.Body.String()); got != "meta,delta,done" {
		t.Fatalf("HTTP SSE 事件顺序 = %q，响应=%s", got, response.Body.String())
	}
	if history.appendCalls != 1 {
		t.Fatalf("HTTP done 前应提交一轮历史，append=%d", history.appendCalls)
	}
}

// newTestCustomerSupportService 创建只用于离线测试的 Fake 依赖，不连接 Redis 或模型服务。
func newTestCustomerSupportService(t *testing.T) (*customerSupportService, *fakeChatProvider, *fakeHistoryStore) {
	t.Helper()
	provider := &fakeChatProvider{}
	history := &fakeHistoryStore{}
	knowledge := businessKnowledge{Contexts: map[customerIntent]string{
		intentProductAdvice:  "商品建议",
		intentDeliveryReturn: "配送退换",
		intentAfterSales:     "售后升级",
		intentGeneral:        "通用回答",
	}}
	service, err := newCustomerSupportService(provider, provider, history, knowledge, governanceConfig{
		Timeout:        time.Second,
		MaxRetries:     0,
		InitialBackoff: time.Millisecond,
		Limiter:        rate.NewLimiter(rate.Inf, 1),
	})
	if err != nil {
		t.Fatalf("创建测试客服服务失败: %v", err)
	}
	return service, provider, history
}

// fakeChatProvider 记录调用次数，以证明未完成骨架不越过模型边界。
type fakeChatProvider struct {
	generateCalls   int
	streamCalls     int
	generateResults []modelResult
	generateErrors  []error
	streamResults   []messageStream
	streamErrors    []error
}

// Generate 记录非流式调用；测试失败时返回可识别错误。
func (provider *fakeChatProvider) Generate(context.Context, modelCall) (modelResult, error) {
	callIndex := provider.generateCalls
	provider.generateCalls++
	if callIndex < len(provider.generateErrors) && provider.generateErrors[callIndex] != nil {
		return modelResult{}, provider.generateErrors[callIndex]
	}
	if callIndex < len(provider.generateResults) {
		return provider.generateResults[callIndex], nil
	}
	return modelResult{}, errors.New("fake provider generate 未配置")
}

// Stream 记录流创建调用；测试失败时返回可识别错误。
func (provider *fakeChatProvider) Stream(context.Context, modelCall) (messageStream, error) {
	callIndex := provider.streamCalls
	provider.streamCalls++
	if callIndex < len(provider.streamErrors) && provider.streamErrors[callIndex] != nil {
		return nil, provider.streamErrors[callIndex]
	}
	if callIndex < len(provider.streamResults) {
		return provider.streamResults[callIndex], nil
	}
	return nil, errors.New("fake provider stream 未配置")
}

// fakeHistoryStore 记录历史写入次数，不实现 Redis 网络行为。
type fakeHistoryStore struct {
	loadCalls            int
	appendCalls          int
	lastLoadUserID       string
	lastLoadSessionID    string
	lastAppendUserID     string
	lastAppendSessionID  string
	lastUserMessage      *schema.Message
	lastAssistantMessage *schema.Message
	history              []*schema.Message
	loadErr              error
	appendErr            error
}

// Load 返回空的 typed 历史，符合首次会话的正常语义。
func (store *fakeHistoryStore) Load(_ context.Context, userID string, sessionID string) ([]*schema.Message, error) {
	store.loadCalls++
	store.lastLoadUserID = userID
	store.lastLoadSessionID = sessionID
	if store.loadErr != nil {
		return nil, store.loadErr
	}
	return store.history, nil
}

// AppendTurn 记录完整轮次写入，用于断言失败路径不会产生副作用。
func (store *fakeHistoryStore) AppendTurn(_ context.Context, userID string, sessionID string, userMessage *schema.Message, assistantMessage *schema.Message) error {
	store.appendCalls++
	store.lastAppendUserID = userID
	store.lastAppendSessionID = sessionID
	store.lastUserMessage = userMessage
	store.lastAssistantMessage = assistantMessage
	return store.appendErr
}

// redisHistoryHook 在不连接真实 Redis 的情况下记录历史读写命令。
type redisHistoryHook struct {
	lrangeValues []string
	processArgs  []any
	pipelineArgs [][]any
}

// DialHook 保留 go-redis 默认连接逻辑；测试命令均由其他 Hook 提前返回。
func (hook *redisHistoryHook) DialHook(next redis.DialHook) redis.DialHook {
	return next
}

// ProcessHook 返回预设的 LRANGE 数据并记录命令参数。
func (hook *redisHistoryHook) ProcessHook(_ redis.ProcessHook) redis.ProcessHook {
	return func(_ context.Context, cmd redis.Cmder) error {
		hook.processArgs = append([]any(nil), cmd.Args()...)
		result, ok := cmd.(*redis.StringSliceCmd)
		if !ok {
			return errors.New("测试 Redis Hook 只支持 StringSliceCmd")
		}
		result.SetVal(append([]string(nil), hook.lrangeValues...))
		return nil
	}
}

// ProcessPipelineHook 记录事务 Pipeline 内的 RPUSH、LTRIM 和 EXPIRE 命令。
func (hook *redisHistoryHook) ProcessPipelineHook(_ redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(_ context.Context, commands []redis.Cmder) error {
		hook.pipelineArgs = make([][]any, 0, len(commands))
		for _, command := range commands {
			hook.pipelineArgs = append(hook.pipelineArgs, append([]any(nil), command.Args()...))
		}
		return nil
	}
}

// recordingEmitter 提供不提交 HTTP 响应的离线事件接收器。
type recordingEmitter struct {
	started    bool
	events     []recordedEvent
	beforeSend func(string)
}

// recordedEvent 保存离线断言所需的 SSE 事件名和原始结构化负载。
type recordedEvent struct {
	name    string
	payload any
}

// Send 标记事件已发送；未完成路径不应调用该方法。
func (emitter *recordingEmitter) Send(event string, payload any) error {
	if emitter.beforeSend != nil {
		emitter.beforeSend(event)
	}
	emitter.started = true
	emitter.events = append(emitter.events, recordedEvent{name: event, payload: payload})
	return nil
}

// Started 返回测试发送器是否已收到任何事件。
func (emitter *recordingEmitter) Started() bool {
	return emitter.started
}

// eventNames 返回记录事件的逗号分隔顺序，便于断言 SSE 协议。
func (emitter *recordingEmitter) eventNames() []string {
	names := make([]string, 0, len(emitter.events))
	for _, event := range emitter.events {
		names = append(names, event.name)
	}
	return names
}

// fakeMessageStream 以预设步骤模拟模型流，不依赖任何网络连接。
type fakeMessageStream struct {
	steps  []streamStep
	index  int
	closed bool
}

// streamStep 表示一次 Recv 的模型块或预设错误。
type streamStep struct {
	chunk modelChunk
	err   error
}

// Recv 按配置顺序返回流步骤，全部消费后模拟正常 EOF。
func (stream *fakeMessageStream) Recv() (modelChunk, error) {
	if stream.index >= len(stream.steps) {
		return modelChunk{}, io.EOF
	}
	step := stream.steps[stream.index]
	stream.index++
	return step.chunk, step.err
}

// Close 标记测试流已关闭，供资源生命周期断言使用。
func (stream *fakeMessageStream) Close() {
	stream.closed = true
}

// temporaryNetworkError 表示允许重试的临时网络失败。
type temporaryNetworkError struct{}

// Error 返回临时网络失败的稳定测试文本。
func (temporaryNetworkError) Error() string {
	return "temporary network failure"
}

// Timeout 表示该测试错误不是尝试级超时。
func (temporaryNetworkError) Timeout() bool {
	return false
}

// Temporary 标记该测试错误属于允许重试的临时网络错误。
func (temporaryNetworkError) Temporary() bool {
	return true
}

// requireCompletedExercise 在当前骨架未完成时跳过未来行为断言，完成 TODO 后自动启用。
func requireCompletedExercise(t *testing.T, err error) {
	t.Helper()
	if errors.Is(err, errExerciseIncomplete) {
		t.Skip("核心 TODO 尚未完成；完成后自动启用该行为断言")
	}
}

// sseEventNames 从 HTTP 响应中提取 SSE event 行的顺序。
func sseEventNames(body string) string {
	names := make([]string, 0)
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(line, "event: ") {
			names = append(names, strings.TrimPrefix(line, "event: "))
		}
	}
	return strings.Join(names, ",")
}

const validClassificationJSON = `{"intent":"delivery_return","response_style":"guided","requires_handoff":false}`
