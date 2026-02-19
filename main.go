package main

import (
	"ecommerce/container"
	"ecommerce/internal/config"
	"ecommerce/internal/logger"
	"ecommerce/internal/mysql"
	"ecommerce/internal/redis"
	"ecommerce/router"
	"ecommerce/util"
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
)

// main 是应用程序入口。
func main() {
	//初始化配置文件
	configErr := config.Init()
	if configErr != nil {
		slog.Error("初始化配置文件错误", "error", configErr)
		panic("初始化配置文件错误")
	}
	// 从配置文件初始化日志
	if err := logger.InitLogConfig(); err != nil {
		slog.Error("日志初始化失败", "error", err)
	}
	log := logger.GetLogger()
	log.Debug("调试信息", "user_id", 123, "action", "login")
	//初始化雪花 ID
	util.Init(config.Cfg.App.MachineId)
	id := util.GenID()
	fmt.Println(id)
	// 初始化Redis客户端
	if err := redis.Init(); err != nil {
		log.Error("Redis初始化失败", "error", err)
		return
	}
	defer func() {
		if err := redis.Close(); err != nil {
		}
	}()

	// 2. 初始化 MySQL（GORM）
	if err := mysql.InitMySQL(); err != nil {
		log.Error("初始化 MySQL 失败", "error", err)
		return
	}
	defer func() {
		if err := mysql.Close(); err != nil {
		}
	}()
	// 3. 初始化容器 依赖注入
	ctn := container.GetInstance()
	// 初始化路由
	r := router.InitTotalRouter(log, ctn)
	// 启动pprof 语言内置的性能分析工具包，用于在运行时收集程序的 CPU、内存、goroutine、阻塞、互斥锁
	if config.Cfg.App.Env != "prod" {
		go func() {
			log.Info("🚀 pprof 启动: http://localhost:6060/debug/pprof/")
			if err := http.ListenAndServe("localhost:6060", nil); err != nil {
				log.Error("pprof 服务启动失败", "error", err)
			}
		}()
	}
	// 启动服务器
	log.Info("服务器启动成功", "APP-NAME", config.Cfg.App.Name, "APP-PORT", config.Cfg.App.Port)
	// 服务器启动是阻塞操作，只有失败才会返回
	if err := r.Run(fmt.Sprintf(":%d", config.Cfg.App.Port)); err != nil {
		log.Error("服务器启动失败", "error", err)
	}

}
