package main

import (
	"context"
	"errors"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// MySQL 与 Qdrant 配置集中放在顶部，练习时直接替换占位值。
const (
	mysqlDSN         = "root:replace-with-your-password@tcp(localhost:3306)/ecommerce?charset=utf8mb4&parseTime=True&loc=Local"
	qdrantBaseURL    = "http://localhost:6333"
	qdrantAPIKey     = "replace-with-your-qdrant-api-key"
	qdrantCollection = "ai_phase3_versions"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type documentVersion struct {
	DocumentID  string
	Version     int
	ContentHash string
	Status      string
	IndexID     string
}

type indexLifecycleConfig struct {
	MySQLDSN         string
	QdrantBaseURL    string
	QdrantAPIKey     string
	QdrantCollection string
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

// validateDocumentVersion 校验状态机、版本和幂等字段。
func validateDocumentVersion(version documentVersion) error {
	// TODO 1：定义状态机、内容哈希、版本号和幂等键契约。
	return errExerciseIncomplete
}

// newMySQLStateDB 使用练习专用 DSN 创建真实 GORM/MySQL 连接；调用方负责关闭底层连接池。
func newMySQLStateDB(dsn string) (*gorm.DB, error) {
	if strings.TrimSpace(dsn) == "" {
		return nil, errors.New("MySQL DSN 不能为空")
	}
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

// newMySQLStateStore 使用真实 GORM 连接创建任务状态存储。
func newMySQLStateStore(db *gorm.DB) (documentStateStore, error) {
	// TODO 2：实现 MySQL Adapter，通过幂等键创建或恢复导入任务。
	return nil, errExerciseIncomplete
}

// writeIsolatedVersion 写入隔离索引，并在失败时记录步骤和原因。
func writeIsolatedVersion(ctx context.Context, states documentStateStore, index vectorIndex, version documentVersion) (string, error) {
	// TODO 3：调用 WriteVersion，失败时通过 MarkFailed 保存可诊断状态。
	return "", errExerciseIncomplete
}

// runImport 执行可恢复的版本化导入，不在索引成功前暴露新版本。
func runImport(ctx context.Context, states documentStateStore, index vectorIndex, version documentVersion) error {
	// TODO 4：索引全部成功后再原子 MarkAvailable，旧版本延迟清理。
	return errExerciseIncomplete
}

// verifyImportRecovery 验证更新、删除、重试和补偿路径。
func verifyImportRecovery(ctx context.Context, config indexLifecycleConfig) error {
	// TODO 5：模拟中断后恢复，并验证重复执行不会暴露半成品版本。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“索引生命周期”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return verifyImportRecovery(ctx, indexLifecycleConfig{
		MySQLDSN:         mysqlDSN,
		QdrantBaseURL:    qdrantBaseURL,
		QdrantAPIKey:     qdrantAPIKey,
		QdrantCollection: qdrantCollection,
	})
}
