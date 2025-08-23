package main

import (
	"context"
	"ecommerce/internal/config"
	"ecommerce/mysql"
	"ecommerce/redis"
	"ecommerce/router"
	"fmt"
	"log"
)

func main() {
	//初始化配置文件
	err1 := config.Init()
	if err1 != nil {
		log.Printf("初始化配置文件错误")
		return
	}
	// 初始化Redis客户端
	if err := redis.Init(); err != nil {
		log.Fatalf("Redis初始化失败: %v", err)
	}
	defer redis.Close() // 程序退出时关闭连接

	// 2. 初始化 MySQL（GORM）
	if err := mysql.InitMySQL(); err != nil {
		log.Fatalf("初始化 MySQL 失败: %v", err)
	}
	defer mysql.Close()

	// 存储字符串
	ctx := context.Background()
	err := redis.Client().Set(ctx, "mykey", "Hello Redis!", 0).Err()
	if err != nil {
		log.Printf("存储到Redis失败: %v", err)
		return
	}
	log.Println("字符串已成功存储到Redis")
	// 初始化路由
	r := router.InitTotalRouter()

	// 启动服务器
	fmt.Println("Server is running on :8080")
	// 服务器启动是阻塞操作，只有失败才会返回
	if err := r.Run(":8080"); err != nil {
		fmt.Printf("Server start failed: %v\n", err)
	}

}
