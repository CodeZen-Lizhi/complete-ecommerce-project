package main

import (
	"context"
	"fmt"
)

// main 启动阶段 3 第 8 个“问题改写、Multi-Query 与 Metadata Filter”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 3 练习 8 未完成: %v\n", err)
	}
}
