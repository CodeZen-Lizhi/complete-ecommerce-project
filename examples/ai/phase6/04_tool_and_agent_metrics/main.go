package main

import (
	"context"
	"fmt"
)

// main 启动阶段 6 第 4 个“Tool 与 Agent 指标”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 6 练习 4 未完成: %v\n", err)
	}
}
