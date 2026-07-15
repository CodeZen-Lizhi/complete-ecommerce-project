package main

import (
	"context"
	"fmt"
)

// main 启动阶段 4 第 6 个“工具失败治理”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 4 练习 6 未完成: %v\n", err)
	}
}
