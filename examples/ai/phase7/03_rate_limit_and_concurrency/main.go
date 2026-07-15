package main

import (
	"context"
	"fmt"
)

// main 启动阶段 7 第 3 个“限流与并发控制”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 7 练习 3 未完成: %v\n", err)
	}
}
