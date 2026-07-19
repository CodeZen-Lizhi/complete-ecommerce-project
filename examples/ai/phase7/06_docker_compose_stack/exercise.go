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

// buildComposeContract 定义服务、卷和环境变量契约。
func buildComposeContract() (composeContract, error) {
	// TODO 1：包含 MySQL、Redis、向量库和应用所需配置。
	return composeContract{}, errExerciseIncomplete
}

// validateInfrastructureServices 校验基础设施服务定义。
func validateInfrastructureServices(path string, contract composeContract) error {
	// TODO 2：检查持久卷、健康检查、最小端口和 Secret 占位符。
	return errExerciseIncomplete
}

// validateApplicationService 校验应用依赖和启动条件。
func validateApplicationService(path string) error {
	// TODO 3：要求 app 使用 depends_on 健康条件等待依赖。
	return errExerciseIncomplete
}

// validateComposeFile 校验服务、健康检查、持久卷和环境变量占位符。
func validateComposeFile(path string, contract composeContract) error {
	// TODO 4：运行 docker compose config，并执行依赖健康烟测。
	return errExerciseIncomplete
}

// composeCommands 返回配置校验、启动、停止和保留数据的安全命令。
func composeCommands(path string) ([]string, error) {
	// TODO 5：返回安全启动/停止命令，清理卷必须使用显式命令。
	return nil, errExerciseIncomplete
}

// runExercise 按执行顺序组织“Docker Compose 本地栈”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	_, err := composeCommands("")
	return err
}
