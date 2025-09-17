package main

import (
	"ecommerce/internal/config"
	"ecommerce/internal/logger"
	"ecommerce/internal/mysql"
	"ecommerce/internal/redis"
	"ecommerce/router"
	"ecommerce/util"
	"fmt"
	"log/slog"
)

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
		log.Error("Redis初始化失败", err)
	}
	defer func() {
		if err := redis.Close(); err != nil {
		}
	}()

	// 2. 初始化 MySQL（GORM）
	if err := mysql.InitMySQL(); err != nil {
		log.Error("初始化 MySQL 失败: %v", err)
	}
	defer func() {
		if err := mysql.Close(); err != nil {
		}
	}()
	// 初始化路由
	r := router.InitTotalRouter()
	// 启动服务器
	log.Info("服务器启动成功", "APP-NAME", config.Cfg.App.Name, "APP-PORT", config.Cfg.App.Port)
	// 服务器启动是阻塞操作，只有失败才会返回
	if err := r.Run(":8080"); err != nil {
		log.Error("服务器启动失败", "error", err)
	}

}
