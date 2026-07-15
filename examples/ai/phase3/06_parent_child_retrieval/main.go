package main

import (
	"context"
	"fmt"
)

// main 启动阶段 3 第 6 个“Parent-Child Retrieval”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 3 练习 6 未完成: %v\n", err)
	}
}
