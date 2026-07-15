package main

import (
	"context"
	"fmt"
)

// main 启动阶段 6 第 6 个“回归对比”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 6 练习 6 未完成: %v\n", err)
	}
}
