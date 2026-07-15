package main

import (
	"context"
	"fmt"
)

// main 启动阶段 3 第 11 个“索引生命周期”练习。
func main() {
	if err := runExercise(context.Background()); err != nil {
		fmt.Printf("阶段 3 练习 11 未完成: %v\n", err)
	}
}
