package main

import (
	"context"
	"errors"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type composeContract struct {
	RequiredServices []string
	RequiredVolumes  []string
	RequiredEnvKeys  []string
}

// validateComposeFile 校验服务、健康检查、持久卷和环境变量占位符。
func validateComposeFile(path string, contract composeContract) error {
	return errExerciseIncomplete
}

// composeCommands 返回配置校验、启动、停止和保留数据的安全命令。
func composeCommands(path string) ([]string, error) {
	return nil, errExerciseIncomplete
}

// runExercise 按执行顺序组织“Docker Compose 本地栈”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：定义 MySQL、Redis、向量库和应用的环境变量契约。
	// TODO 2：补齐各服务持久卷、健康检查和最小端口。
	// TODO 3：添加 app 服务，并用 depends_on 健康条件等待依赖。
	// TODO 4：实现 validateComposeFile，并运行 docker compose config 与健康烟测。
	// TODO 5：实现 composeCommands，停止默认保留数据，清理卷需显式命令。
	return errExerciseIncomplete
}
