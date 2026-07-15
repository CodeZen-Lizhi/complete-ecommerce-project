package main

import (
	"context"
	"fmt"
)

// main 启动阶段 4 第 5 个“多工具注册与连续调用”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 4 练习 5 未完成: %v\n", err)
	}
}
