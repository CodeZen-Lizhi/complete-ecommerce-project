package main

import (
	"context"
	"fmt"
)

// main 启动阶段 3 第 4 个“向量索引与检索”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 3 练习 4 未完成: %v\n", err)
	}
}
