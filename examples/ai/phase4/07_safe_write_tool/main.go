package main

import (
	"context"
	"fmt"
)

// main 启动阶段 4 第 7 个“安全写操作模拟”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 4 练习 7 未完成: %v\n", err)
	}
}
