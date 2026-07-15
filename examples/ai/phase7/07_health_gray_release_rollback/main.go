package main

import (
	"context"
	"fmt"
)

// main 启动阶段 7 第 7 个“健康检查、灰度与回滚”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 7 练习 7 未完成: %v\n", err)
	}
}
