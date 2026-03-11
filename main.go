package main

import (
	"context"
	"ecommerce/container"
	"ecommerce/internal/config"
	"ecommerce/internal/logger"
	"ecommerce/internal/mysql"
	"ecommerce/internal/redis"
	"ecommerce/router"
	"ecommerce/util"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// main 是应用程序入口
func main() {
	// 初始化配置文件
	configErr := config.Init()
	if configErr != nil {
		slog.Error("初始化配置文件错误", "error", configErr)
		panic("初始化配置文件错误")
	}

	// 初始化日志
	if err := logger.InitLogConfig(); err != nil {
		slog.Error("日志初始化失败", "error", err)
	}
	log := logger.GetLogger()

	// 初始化雪花 ID 生成器
	util.Init(config.Cfg.App.MachineId)

	// 初始化 Redis 客户端
	if err := redis.Init(); err != nil {
		log.Error("Redis 初始化失败", "error", err)
		return
	}
	defer func() {
		if err := redis.Close(); err != nil {
			log.Error("Redis 关闭失败", "error", err)
		}
	}()

	// 初始化 MySQL 数据库
	if err := mysql.InitMySQL(); err != nil {
		log.Error("MySQL 初始化失败", "error", err)
		return
	}
	defer func() {
		if err := mysql.Close(); err != nil {
			log.Error("MySQL 关闭失败", "error", err)
		}
	}()

	// 初始化依赖注入容器
	ctn := container.GetInstance(mysql.DB)

	// 初始化路由
	r := router.InitTotalRouter(log, ctn)

	// 启动 pprof 性能分析工具（非生产环境）
	if config.Cfg.App.Env != "prod" {
		go func() {
			log.Info("pprof 已启动", "url", "http://localhost:6060/debug/pprof/")
			if err := http.ListenAndServe("localhost:6060", nil); err != nil {
				log.Error("pprof 服务启动失败", "error", err)
			}
		}()
	}

	// 创建 HTTP Server
	addr := fmt.Sprintf(":%d", config.Cfg.App.Port)
	server := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// 监听退出信号（Ctrl+C / SIGTERM）
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 启动 HTTP 服务器
	go func() {
		log.Info("服务器启动成功", "name", config.Cfg.App.Name, "port", config.Cfg.App.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("服务器启动失败", "error", err)
		}
	}()

	// 等待退出信号
	<-ctx.Done()
	log.Info("收到退出信号，开始优雅停机")

	// 优雅停机
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("服务器优雅停机失败", "error", err)
		return
	}
	log.Info("服务器已优雅停机")
}
