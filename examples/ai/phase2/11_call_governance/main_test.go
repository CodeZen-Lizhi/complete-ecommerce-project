package main

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"golang.org/x/time/rate"
)

type scriptedResult struct {
	response callResponse
	err      error
}

type scriptedModel struct {
	results      []scriptedResult
	calls        int
	hadDeadlines []bool
	contexts     []context.Context
}

// Generate 按顺序返回预设结果，并记录每次尝试是否携带截止时间。
func (m *scriptedModel) Generate(ctx context.Context, _ callRequest) (callResponse, error) {
	m.calls++
	_, hasDeadline := ctx.Deadline()
	m.hadDeadlines = append(m.hadDeadlines, hasDeadline)
	m.contexts = append(m.contexts, ctx)
	if m.calls > len(m.results) {
		return callResponse{}, errors.New("脚本化模型缺少预设结果")
	}
	result := m.results[m.calls-1]
	return result.response, result.err
}

// testGovernanceConfig 返回不限制测试速度的治理配置。
func testGovernanceConfig(maxRetries int) governanceConfig {
	return governanceConfig{
		Timeout:        100 * time.Millisecond,
		MaxRetries:     maxRetries,
		InitialBackoff: time.Millisecond,
		Limiter:        rate.NewLimiter(rate.Inf, 1),
	}
}

// TestGovernedCallSuccess 验证首轮成功时的响应、指标和尝试超时。
func TestGovernedCallSuccess(t *testing.T) {
	client := &scriptedModel{results: []scriptedResult{{
		response: callResponse{Text: "ok", PromptTokens: 2, OutputTokens: 3},
	}}}

	response, metrics, err := governedCall(
		context.Background(),
		client,
		testGovernanceConfig(2),
		callRequest{Prompt: "test"},
	)
	if err != nil {
		t.Fatalf("governedCall() error = %v", err)
	}
	if response.Text != "ok" {
		t.Fatalf("governedCall() Text = %q，期望 ok", response.Text)
	}
	if metrics.Attempts != 1 || !metrics.Success || metrics.Tokens != 5 {
		t.Fatalf("governedCall() metrics = %+v", metrics)
	}
	if len(client.hadDeadlines) != 1 || !client.hadDeadlines[0] {
		t.Fatalf("Generate() deadline 记录 = %v，期望首轮携带截止时间", client.hadDeadlines)
	}
	select {
	case <-client.contexts[0].Done():
	default:
		t.Fatal("首轮 Generate() Context 未在本轮结束后取消")
	}
}

// TestGovernedCallRetryThenSuccess 验证 429 后退避一次并成功返回累计指标。
func TestGovernedCallRetryThenSuccess(t *testing.T) {
	client := &scriptedModel{results: []scriptedResult{
		{err: &openai.APIError{HTTPStatusCode: http.StatusTooManyRequests}},
		{response: callResponse{Text: "ok", PromptTokens: 4, OutputTokens: 6}},
	}}

	response, metrics, err := governedCall(
		context.Background(),
		client,
		testGovernanceConfig(1),
		callRequest{Prompt: "test"},
	)
	if err != nil {
		t.Fatalf("governedCall() error = %v", err)
	}
	if response.Text != "ok" || metrics.Attempts != 2 || !metrics.Success || metrics.Tokens != 10 {
		t.Fatalf("governedCall() response = %+v, metrics = %+v", response, metrics)
	}
	for attempt, hasDeadline := range client.hadDeadlines {
		if !hasDeadline {
			t.Errorf("第 %d 次 Generate() 未携带截止时间", attempt+1)
		}
	}
	for attempt, attemptCtx := range client.contexts {
		select {
		case <-attemptCtx.Done():
		default:
			t.Errorf("第 %d 次 Generate() Context 未在本轮结束后取消", attempt+1)
		}
	}
}

// TestGovernedCallExhaustionKeepsLastError 验证耗尽重试时保留最后错误和失败指标。
func TestGovernedCallExhaustionKeepsLastError(t *testing.T) {
	firstErr := &openai.APIError{HTTPStatusCode: http.StatusServiceUnavailable, Message: "first"}
	lastErr := &openai.APIError{HTTPStatusCode: http.StatusBadGateway, Message: "last"}
	client := &scriptedModel{results: []scriptedResult{{err: firstErr}, {err: lastErr}}}

	_, metrics, err := governedCall(
		context.Background(),
		client,
		testGovernanceConfig(1),
		callRequest{Prompt: "test"},
	)
	if err == nil {
		t.Fatal("governedCall() error = nil，期望重试耗尽错误")
	}
	var gotAPIError *openai.APIError
	if !errors.As(err, &gotAPIError) || gotAPIError != lastErr {
		t.Fatalf("governedCall() error = %v，未保留最后错误", err)
	}
	if metrics.Attempts != 2 || metrics.Success {
		t.Fatalf("governedCall() metrics = %+v", metrics)
	}
}

// TestGovernedCallNonRetryableFailure 验证未知错误不会重试并返回失败指标。
func TestGovernedCallNonRetryableFailure(t *testing.T) {
	nonRetryableErr := errors.New("解析错误")
	client := &scriptedModel{results: []scriptedResult{{err: nonRetryableErr}}}

	_, metrics, err := governedCall(
		context.Background(),
		client,
		testGovernanceConfig(2),
		callRequest{Prompt: "test"},
	)
	if !errors.Is(err, nonRetryableErr) {
		t.Fatalf("governedCall() error = %v，期望包含解析错误", err)
	}
	if client.calls != 1 || metrics.Attempts != 1 || metrics.Success {
		t.Fatalf("Generate() calls = %d, metrics = %+v", client.calls, metrics)
	}
}

// TestGovernedCallAttemptTimeoutDoesNotRetry 验证单次模型超时不会触发重复计费风险。
func TestGovernedCallAttemptTimeoutDoesNotRetry(t *testing.T) {
	client := &scriptedModel{results: []scriptedResult{{err: context.DeadlineExceeded}}}

	_, metrics, err := governedCall(
		context.Background(),
		client,
		testGovernanceConfig(2),
		callRequest{Prompt: "test"},
	)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("governedCall() error = %v，期望 context.DeadlineExceeded", err)
	}
	if client.calls != 1 || metrics.Attempts != 1 || metrics.Success {
		t.Fatalf("Generate() calls = %d, metrics = %+v", client.calls, metrics)
	}
}

// TestWaitBackoffCancellation 验证退避等待能够立即响应 Context 取消。
func TestWaitBackoffCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	startedAt := time.Now()
	err := waitBackoff(ctx, time.Second, 0)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("waitBackoff() error = %v，期望 context.Canceled", err)
	}
	if elapsed := time.Since(startedAt); elapsed >= 100*time.Millisecond {
		t.Fatalf("waitBackoff() 取消耗时 = %s，期望立即返回", elapsed)
	}
}

// TestGovernedCallLimiterCancellation 验证等待限流许可时取消不会调用模型。
func TestGovernedCallLimiterCancellation(t *testing.T) {
	limiter := rate.NewLimiter(rate.Every(time.Second), 1)
	if !limiter.Allow() {
		t.Fatal("测试限流器未能消耗初始令牌")
	}
	config := testGovernanceConfig(0)
	config.Limiter = limiter
	client := &scriptedModel{results: []scriptedResult{{response: callResponse{Text: "unexpected"}}}}
	ctx, cancel := context.WithCancel(context.Background())
	cancelTimer := time.AfterFunc(10*time.Millisecond, cancel)
	defer cancelTimer.Stop()

	_, metrics, err := governedCall(ctx, client, config, callRequest{Prompt: "test"})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("governedCall() error = %v，期望 context.Canceled", err)
	}
	if client.calls != 0 || metrics.Attempts != 0 || metrics.Success {
		t.Fatalf("Generate() calls = %d, metrics = %+v", client.calls, metrics)
	}
}
