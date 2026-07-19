package main

import (
	"context"
	"errors"
)

var errExerciseIncomplete = errors.New("练习尚未完成，请按 TODO 顺序实现")

type loadedDocument struct {
	ID       string
	Title    string
	Page     int
	Source   string
	Content  string
	Metadata map[string]string
}

type documentLoader interface {
	// Load 从受控路径加载文档，并返回页级或结构级内容。
	Load(ctx context.Context, path string) ([]loadedDocument, error)
}

type pdfLoaderConfig struct {
	AllowedRoot string
	MaxBytes    int64
}

// validateDocumentPath 校验文件扩展名、大小和允许目录，防止路径穿越。
func validateDocumentPath(path string, allowedRoot string, maxBytes int64) error {
	// TODO 1：定义允许目录、文件类型和大小上限，并拒绝路径穿越。
	return errExerciseIncomplete
}

// newMarkdownLoader 创建保留标题层级和来源的 Markdown Loader。
func newMarkdownLoader(config pdfLoaderConfig) (documentLoader, error) {
	// TODO 2：实现 Markdown Loader，保留标题层级和来源。
	return nil, errExerciseIncomplete
}

// newPDFPageLoader 使用真实 PDF 解析库创建页级 Loader；固定内存内容不能作为实现。
func newPDFPageLoader(config pdfLoaderConfig) (documentLoader, error) {
	// TODO 3：读取本地 fixture 的页级文本，并为解析错误补充页码上下文。
	return nil, errExerciseIncomplete
}

// normalizeDocument 补齐来源、页码、租户和权限 Metadata。
func normalizeDocument(document loadedDocument) (loadedDocument, error) {
	// TODO 4：强制补齐租户、权限、来源和页码 Metadata。
	return loadedDocument{}, errExerciseIncomplete
}

// verifyDocumentFixtures 验证正常文件和加载失败边界。
func verifyDocumentFixtures(ctx context.Context, config pdfLoaderConfig) error {
	// TODO 5：覆盖空文档、路径越界、不支持格式和缺失权限字段。
	return errExerciseIncomplete
}

// runExercise 按执行顺序组织“Markdown/PDF 加载与 Metadata”练习的核心步骤。
func runExercise(ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context 不能为空")
	}

	return verifyDocumentFixtures(ctx, pdfLoaderConfig{})
}
