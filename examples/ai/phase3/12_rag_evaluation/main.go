package main

import (
	"context"
	"fmt"
)

// main 启动阶段 3 第 12 个“RAG 检索评估”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 3 练习 12 未完成: %v\n", err)
	}
}
