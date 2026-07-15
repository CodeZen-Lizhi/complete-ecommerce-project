package main

import (
	"context"
	"fmt"
)

// main 启动阶段 7 第 5 个“Metrics 与告警”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 7 练习 5 未完成: %v\n", err)
	}
}
