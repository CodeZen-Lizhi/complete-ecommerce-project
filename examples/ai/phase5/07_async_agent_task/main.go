package main

import (
	"context"
	"fmt"
)

// main 启动阶段 5 第 7 个“异步 Agent 任务”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 5 练习 7 未完成: %v\n", err)
	}
}
