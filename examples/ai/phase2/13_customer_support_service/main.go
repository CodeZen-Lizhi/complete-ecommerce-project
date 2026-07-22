package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const (
	redisPingTimeout    = 5 * time.Second
	shutdownGracePeriod = 10 * time.Second
)

// main 建立可取消的进程 Context，并启动独立客服练习服务。
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := runExercise(ctx); err != nil {
		slog.Error("客服综合练习退出", "error", err)
	}
}

// runExercise 按配置、Redis、知识、Provider、路由的顺序装配服务，任一失败立即返回。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return fmt.Errorf("练习 Context 不能为空")
	}
	config, err := loadExerciseConfig(os.Getenv)
	if err != nil {
		return err
	}
	redisClient := redis.NewClient(&redis.Options{Addr: config.RedisAddr})
	defer func() {
		if closeErr := redisClient.Close(); closeErr != nil {
			slog.Warn("关闭练习 Redis Client 失败", "error", closeErr)
		}
	}()
	if err := pingRedis(ctx, redisClient, config.CallTimeout); err != nil {
		return err
	}
	knowledge := defaultBusinessKnowledge()
	classificationFormat, err := buildClassificationResponseFormat()
	if err != nil {
		return err
	}
	classificationProvider, err := newOpenAIProvider(ctx, providerConfig{
		BaseURL:        config.ModelBaseURL,
		APIKey:         config.ModelAPIKey,
		Model:          config.ModelName,
		Timeout:        config.CallTimeout,
		ResponseFormat: classificationFormat,
	})
	if err != nil {
		return err
	}
	answerProvider, err := newOpenAIProvider(ctx, providerConfig{
		BaseURL: config.ModelBaseURL,
		APIKey:  config.ModelAPIKey,
		Model:   config.ModelName,
		Timeout: config.CallTimeout,
	})
	if err != nil {
		return err
	}
	history, err := newRedisHistoryStore(redisClient, config.HistoryTurns, config.HistoryTTL)
	if err != nil {
		return err
	}
	service, err := newCustomerSupportService(
		classificationProvider,
		answerProvider,
		history,
		knowledge,
		newGovernanceConfig(config),
	)
	if err != nil {
		return err
	}
	return serveHTTP(ctx, config.ListenAddr, newCustomerSupportRouter(service))
}

// pingRedis 在有界 Context 内确认 Redis 可用，避免服务在缺失会话存储时启动。
func pingRedis(ctx context.Context, client *redis.Client, callTimeout time.Duration) error {
	if ctx == nil || client == nil {
		return fmt.Errorf("Redis Ping 依赖不能为空")
	}
	timeout := redisPingTimeout
	if callTimeout > 0 && callTimeout < timeout {
		timeout = callTimeout
	}
	pingContext, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	if err := client.Ping(pingContext).Err(); err != nil {
		return fmt.Errorf("Redis Ping 失败: %w", err)
	}
	return nil
}

// serveHTTP 仅监听已校验的回环地址，并在 Context 取消后有界优雅关闭。
func serveHTTP(ctx context.Context, address string, router *gin.Engine) error {
	if ctx == nil || router == nil {
		return fmt.Errorf("HTTP 服务依赖不能为空")
	}
	server := &http.Server{
		Addr:              address,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}
	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("启动 HTTP 服务失败: %w", err)
	case <-ctx.Done():
		shutdownContext, cancel := context.WithTimeout(context.Background(), shutdownGracePeriod)
		defer cancel()
		if err := server.Shutdown(shutdownContext); err != nil {
			return fmt.Errorf("关闭 HTTP 服务失败: %w", err)
		}
		if err := <-serverErrors; err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("HTTP 服务异常退出: %w", err)
		}
		return nil
	}
}
