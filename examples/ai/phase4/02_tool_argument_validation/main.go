package main

import (
	"context"
	"fmt"
)

// main 启动阶段 4 第 2 个“Tool 参数严格校验”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 4 练习 2 未完成: %v\n", err)
	}
}
