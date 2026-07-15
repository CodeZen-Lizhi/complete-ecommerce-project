# 阶段 3 练习 2：Markdown/PDF 加载与 Metadata

## 练习目标

加载 Markdown/PDF，保留标题、页码、来源、租户和权限 Metadata。

## 前置知识

- Eino 组件的 Context 与错误传播。
- 文档、Chunk、Metadata、向量和 Top-K 的基本含义。
- 默认配置只用于本地骨架验证，真实模型或向量库必须显式配置。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：定义支持的文件类型、大小上限和允许目录。
2. TODO 2：实现 Markdown 读取并提取标题层级。
3. TODO 3：通过可替换 Loader 接口接入 PDF 页级文本。
4. TODO 4：规范化来源、页码、租户和权限 Metadata。
5. TODO 5：拒绝空文档、越界路径和缺失权限字段，并补 fixture 测试。

## 开始练习

```bash
go run ./examples/ai/phase3/02_document_loading_metadata
```

骨架默认返回“练习尚未完成”，不会连接模型、数据库或向量库。

## 验证方式

```bash
gofmt -w examples/ai/phase3/02_document_loading_metadata/*.go
go test -timeout=60s ./examples/ai/phase3/02_document_loading_metadata
go vet ./examples/ai/phase3/02_document_loading_metadata
```

## 完成标准

- 输入、空结果、维度、Top-K、权限和取消边界得到明确处理。
- 外部依赖错误使用 `%w` 保留原因，不吞掉部分失败。
- 结果可由固定 fixture 重复验证，并记录必要质量或延迟指标。
- 不提交真实 API Key、数据库密码或向量库凭证。

## 暂不实现

- OCR、复杂表格恢复和云对象存储。
