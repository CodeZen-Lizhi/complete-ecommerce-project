package main

import (
	"context"
	"fmt"
)

// main 启动阶段 3 第 9 个“候选 Rerank”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 3 练习 9 未完成: %v\n", err)
	}
}
