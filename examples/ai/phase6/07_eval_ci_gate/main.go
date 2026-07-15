package main

import (
	"context"
	"fmt"
)

// main 启动阶段 6 第 7 个“CI 质量门禁”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 6 练习 7 未完成: %v\n", err)
	}
}
