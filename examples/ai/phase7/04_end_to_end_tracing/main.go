package main

import (
	"context"
	"fmt"
)

// main 启动阶段 7 第 4 个“端到端 Trace”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 7 练习 4 未完成: %v\n", err)
	}
}
