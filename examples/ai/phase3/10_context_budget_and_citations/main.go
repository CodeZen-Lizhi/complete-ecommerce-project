package main

import (
	"context"
	"fmt"
)

// main 启动阶段 3 第 10 个“上下文预算与引用”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 3 练习 10 未完成: %v\n", err)
	}
}
