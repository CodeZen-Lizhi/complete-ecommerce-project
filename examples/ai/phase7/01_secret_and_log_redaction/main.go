package main

import (
	"context"
	"fmt"
)

// main 启动阶段 7 第 1 个“Secret 管理与日志脱敏”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 7 练习 1 未完成: %v\n", err)
	}
}
