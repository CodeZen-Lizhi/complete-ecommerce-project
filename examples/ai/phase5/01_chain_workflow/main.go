package main

import (
	"context"
	"fmt"
)

// main 启动阶段 5 第 1 个“固定顺序 Chain”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 5 练习 1 未完成: %v\n", err)
	}
}
