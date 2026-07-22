package main

import (
	"fmt"
	"math"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	defaultListenAddr       = "127.0.0.1:8093"
	defaultModelBaseURL     = "http://localhost:8084/v1"
	defaultModelName        = "gpt-5.4-mini"
	defaultRedisAddr        = "127.0.0.1:6379"
	defaultHistoryTTL       = 30 * time.Minute
	defaultHistoryTurns     = 6
	defaultCallTimeout      = 30 * time.Second
	defaultMaxRetries       = 2
	defaultInitialBackoff   = 200 * time.Millisecond
	defaultCallsPerSecond   = 2.0
	maxMessageLength        = 4000
	maximumIdentifierLength = 128
)

// exerciseConfig 保存独立客服练习的环境配置；敏感 API Key 只来自运行环境。
type exerciseConfig struct {
	ListenAddr     string
	ModelBaseURL   string
	ModelAPIKey    string
	ModelName      string
	RedisAddr      string
	HistoryTTL     time.Duration
	HistoryTurns   int
	CallTimeout    time.Duration
	MaxRetries     int
	InitialBackoff time.Duration
	CallsPerSecond float64
}

// loadExerciseConfig 从环境变量加载配置，并在创建 Redis 或模型客户端前拒绝非法值。
func loadExerciseConfig(getenv func(string) string) (exerciseConfig, error) {
	if getenv == nil {
		return exerciseConfig{}, fmt.Errorf("环境变量读取函数不能为空")
	}

	historyTTL, err := parsePositiveDuration("AI_DEMO_HISTORY_TTL", valueOrDefault(getenv, "AI_DEMO_HISTORY_TTL", defaultHistoryTTL.String()))
	if err != nil {
		return exerciseConfig{}, err
	}
	callTimeout, err := parsePositiveDuration("AI_DEMO_CALL_TIMEOUT", valueOrDefault(getenv, "AI_DEMO_CALL_TIMEOUT", defaultCallTimeout.String()))
	if err != nil {
		return exerciseConfig{}, err
	}
	initialBackoff, err := parsePositiveDuration("AI_DEMO_INITIAL_BACKOFF", valueOrDefault(getenv, "AI_DEMO_INITIAL_BACKOFF", defaultInitialBackoff.String()))
	if err != nil {
		return exerciseConfig{}, err
	}
	historyTurns, err := parsePositiveInt("AI_DEMO_HISTORY_TURNS", valueOrDefault(getenv, "AI_DEMO_HISTORY_TURNS", strconv.Itoa(defaultHistoryTurns)))
	if err != nil {
		return exerciseConfig{}, err
	}
	maxRetries, err := parseNonNegativeInt("AI_DEMO_MAX_RETRIES", valueOrDefault(getenv, "AI_DEMO_MAX_RETRIES", strconv.Itoa(defaultMaxRetries)))
	if err != nil {
		return exerciseConfig{}, err
	}
	callsPerSecond, err := parsePositiveFloat("AI_DEMO_CALLS_PER_SECOND", valueOrDefault(getenv, "AI_DEMO_CALLS_PER_SECOND", strconv.FormatFloat(defaultCallsPerSecond, 'f', -1, 64)))
	if err != nil {
		return exerciseConfig{}, err
	}

	config := exerciseConfig{
		ListenAddr:     valueOrDefault(getenv, "AI_DEMO_LISTEN_ADDR", defaultListenAddr),
		ModelBaseURL:   valueOrDefault(getenv, "AI_DEMO_MODEL_BASE_URL", defaultModelBaseURL),
		ModelAPIKey:    strings.TrimSpace(getenv("OPENAI_API_KEY")),
		ModelName:      valueOrDefault(getenv, "AI_DEMO_MODEL", defaultModelName),
		RedisAddr:      valueOrDefault(getenv, "AI_DEMO_REDIS_ADDR", defaultRedisAddr),
		HistoryTTL:     historyTTL,
		HistoryTurns:   historyTurns,
		CallTimeout:    callTimeout,
		MaxRetries:     maxRetries,
		InitialBackoff: initialBackoff,
		CallsPerSecond: callsPerSecond,
	}
	if err := validateExerciseConfig(config); err != nil {
		return exerciseConfig{}, err
	}
	return config, nil
}

// valueOrDefault 返回去除空白后的环境值；空值使用调用方提供的安全默认值。
func valueOrDefault(getenv func(string) string, key string, fallback string) string {
	if value := strings.TrimSpace(getenv(key)); value != "" {
		return value
	}
	return fallback
}

// parsePositiveDuration 解析必须大于零的 Go Duration 配置。
func parsePositiveDuration(key string, value string) (time.Duration, error) {
	duration, err := time.ParseDuration(strings.TrimSpace(value))
	if err != nil || duration <= 0 {
		return 0, fmt.Errorf("%s 必须是大于零的 Duration", key)
	}
	return duration, nil
}

// parsePositiveInt 解析必须大于零的整型配置。
func parsePositiveInt(key string, value string) (int, error) {
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || parsed <= 0 {
		return 0, fmt.Errorf("%s 必须是大于零的整数", key)
	}
	return parsed, nil
}

// parseNonNegativeInt 解析允许为零的整型配置。
func parseNonNegativeInt(key string, value string) (int, error) {
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || parsed < 0 {
		return 0, fmt.Errorf("%s 必须是非负整数", key)
	}
	return parsed, nil
}

// parsePositiveFloat 解析必须大于零的浮点型配置。
func parsePositiveFloat(key string, value string) (float64, error) {
	parsed, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil || math.IsNaN(parsed) || math.IsInf(parsed, 0) || parsed <= 0 {
		return 0, fmt.Errorf("%s 必须是大于零的数字", key)
	}
	return parsed, nil
}

// validateExerciseConfig 校验地址、密钥和数值边界，保证失败发生在外部客户端创建之前。
func validateExerciseConfig(config exerciseConfig) error {
	if err := validateLoopbackAddress(config.ListenAddr); err != nil {
		return err
	}
	if strings.TrimSpace(config.ModelAPIKey) == "" {
		return fmt.Errorf("OPENAI_API_KEY 未配置")
	}
	modelURL, err := url.ParseRequestURI(config.ModelBaseURL)
	if err != nil || modelURL.Scheme == "" || modelURL.Host == "" {
		return fmt.Errorf("AI_DEMO_MODEL_BASE_URL 不是有效的绝对 URL")
	}
	if modelURL.Scheme != "http" && modelURL.Scheme != "https" {
		return fmt.Errorf("AI_DEMO_MODEL_BASE_URL 只支持 http 或 https")
	}
	if strings.TrimSpace(config.ModelName) == "" {
		return fmt.Errorf("AI_DEMO_MODEL 未配置")
	}
	if _, _, err := net.SplitHostPort(config.RedisAddr); err != nil {
		return fmt.Errorf("AI_DEMO_REDIS_ADDR 必须为 host:port: %w", err)
	}
	if config.HistoryTTL <= 0 || config.HistoryTurns <= 0 || config.CallTimeout <= 0 || config.InitialBackoff <= 0 || config.CallsPerSecond <= 0 || config.MaxRetries < 0 {
		return fmt.Errorf("练习配置包含非法的数值边界")
	}
	return nil
}

// validateLoopbackAddress 确保练习 HTTP 服务只绑定数值形式的回环地址。
func validateLoopbackAddress(address string) error {
	host, _, err := net.SplitHostPort(strings.TrimSpace(address))
	if err != nil {
		return fmt.Errorf("AI_DEMO_LISTEN_ADDR 必须为 host:port: %w", err)
	}
	ip := net.ParseIP(host)
	if ip == nil || !ip.IsLoopback() {
		return fmt.Errorf("AI_DEMO_LISTEN_ADDR 必须绑定回环地址")
	}
	return nil
}
