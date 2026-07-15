package main

import (
	"context"
	"fmt"
)

// main 启动阶段 5 第 4 个“Agent 状态与预算”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 5 练习 4 未完成: %v\n", err)
	}
}
