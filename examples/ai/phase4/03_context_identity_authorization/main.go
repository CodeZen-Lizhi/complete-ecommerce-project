package main

import (
	"context"
	"fmt"
)

// main 启动阶段 4 第 3 个“Context 身份与再次授权”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 4 练习 3 未完成: %v\n", err)
	}
}
