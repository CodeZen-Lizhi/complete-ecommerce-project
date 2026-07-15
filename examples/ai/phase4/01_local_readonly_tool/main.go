package main

import (
	"context"
	"fmt"
)

// main 启动阶段 4 第 1 个“本地只读工具与 ToolsNode”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 4 练习 1 未完成: %v\n", err)
	}
}
