package main

import (
	"context"
	"fmt"
)

// main 启动阶段 6 第 3 个“回答事实与引用检查”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 6 练习 3 未完成: %v\n", err)
	}
}
