package main

import (
	"context"
	"errors"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type documentVersion struct {
	DocumentID  string
	Version     int
	ContentHash string
	Status      string
	IndexID     string
}

type documentStateStore interface {
	// CreateImportTask 按幂等键创建或返回已有导入任务。
	CreateImportTask(ctx context.Context, idempotencyKey string, version documentVersion) (documentVersion, error)
	// MarkAvailable 原子切换可用版本。
	MarkAvailable(ctx context.Context, documentID string, version int, indexID string) error
	// MarkFailed 记录失败步骤、原因和重试次数。
	MarkFailed(ctx context.Context, documentID string, version int, step string, cause error) error
}

type vectorIndex interface {
	// WriteVersion 写入隔离的新版本索引。
	WriteVersion(ctx context.Context, version documentVersion) (string, error)
	// DeleteVersion 删除指定索引版本。
	DeleteVersion(ctx context.Context, indexID string) error
}

// runImport 执行可恢复的版本化导入，不在索引成功前暴露新版本。
func runImport(ctx context.Context, states documentStateStore, index vectorIndex, version documentVersion) error {
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“索引生命周期”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	// TODO 1：定义状态机、内容哈希、版本号和幂等键契约。
	// TODO 2：通过 documentStateStore 创建或恢复导入任务。
	// TODO 3：使用 vectorIndex 写入隔离版本，失败时记录步骤和原因。
	// TODO 4：索引全部成功后再原子 MarkAvailable，旧版本延迟清理。
	// TODO 5：实现更新、删除、重试和补偿测试，验证中断后可恢复。
	return errExerciseIncomplete
}
