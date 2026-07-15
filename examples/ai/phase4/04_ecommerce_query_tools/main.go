package main

import (
	"context"
	"fmt"
)

// main 启动阶段 4 第 4 个“电商业务查询工具”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 4 练习 4 未完成: %v\n", err)
	}
}
