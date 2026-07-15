package main

import (
	"context"
	"fmt"
)

// main 启动阶段 3 第 1 个“Embedding、余弦相似度与 Top-K”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 3 练习 1 未完成: %v\n", err)
	}
}
