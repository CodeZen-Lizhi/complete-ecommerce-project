package main

import (
	"context"
	"fmt"
)

// main 启动阶段 5 第 8 个“任务恢复与幂等”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 5 练习 8 未完成: %v\n", err)
	}
}
