package main

import (
	"context"
	"fmt"
)

// main 启动阶段 3 第 5 个“HNSW 参数调优”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 3 练习 5 未完成: %v\n", err)
	}
}
