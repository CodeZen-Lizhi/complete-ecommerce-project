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

	"github.com/cloudwego/eino/schema"
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
