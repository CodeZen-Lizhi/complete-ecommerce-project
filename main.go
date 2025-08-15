package main

import (
	"ecommerce/router"
	"fmt"
)

func main() {
	// 初始化路由
	r := router.InitTotalRouter()

	// 启动服务器
	fmt.Println("Server is running on :8080")
	// 服务器启动是阻塞操作，只有失败才会返回
	if err := r.Run(":8080"); err != nil {
		fmt.Printf("Server start failed: %v\n", err)
	} else {
		fmt.Println("Server started successfully.")
	}
}
