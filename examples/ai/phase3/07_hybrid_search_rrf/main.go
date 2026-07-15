package main

import (
	"context"
	"fmt"
)

// main 启动阶段 3 第 7 个“Dense + BM25 + RRF”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 3 练习 7 未完成: %v\n", err)
	}
}
