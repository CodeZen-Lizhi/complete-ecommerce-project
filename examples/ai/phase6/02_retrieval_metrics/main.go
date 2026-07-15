package main

import (
	"context"
	"fmt"
)

// main 启动阶段 6 第 2 个“检索指标”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 6 练习 2 未完成: %v\n", err)
	}
}
