package main

import (
	"context"
	"fmt"
)

// main 启动阶段 7 第 2 个“Prompt Injection 防御”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 7 练习 2 未完成: %v\n", err)
	}
}
