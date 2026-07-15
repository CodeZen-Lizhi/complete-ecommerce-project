package main

import (
	"context"
	"fmt"
)

// main 启动阶段 3 第 3 个“切块策略对比”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 3 练习 3 未完成: %v\n", err)
	}
}
