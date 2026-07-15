package main

import (
	"context"
	"fmt"
)

// main 启动阶段 7 第 8 个“依赖故障与降级”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 7 练习 8 未完成: %v\n", err)
	}
}
