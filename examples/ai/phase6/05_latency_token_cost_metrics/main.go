package main

import (
	"context"
	"fmt"
)

// main 启动阶段 6 第 5 个“延迟、Token 与成本指标”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 6 练习 5 未完成: %v\n", err)
	}
}
