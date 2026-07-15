package main

import (
	"context"
	"fmt"
)

// main 启动阶段 7 第 6 个“Docker Compose 本地栈”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 7 练习 6 未完成: %v\n", err)
	}
}
