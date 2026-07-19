# 阶段 3 练习 4：向量索引与检索

## 练习目标

使用 Eino Indexer 写入一种向量存储，并通过 Retriever 完成 Top-K 检索。

## 前置知识

- Eino 组件的 Context 与错误传播。
- 文档、Chunk、Metadata、向量和 Top-K 的基本含义。
- 默认配置只用于本地骨架验证，真实模型或向量库必须显式配置。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：定义文档 ID、Chunk ID、向量维度和 Metadata 契约。
2. TODO 2：创建 Embedder、Indexer 和 Retriever，校验配置兼容性。
3. TODO 3：批量写入 Chunk，处理部分失败和重复 ID。
4. TODO 4：对查询生成向量并执行 Top-K 检索。
5. TODO 5：验证返回内容、分数和 Metadata，并记录空结果。

## 开始练习

先修改 `exercise.go` 顶部的 `embeddingBaseURL`、`embeddingAPIKey`、`embeddingModelName`、`qdrantBaseURL`、`qdrantAPIKey` 和 `qdrantCollection`。

```bash
go run ./examples/ai/phase3/04_vector_index_and_retrieval
```

骨架默认返回“练习尚未完成”，不会连接模型、数据库或向量库。

## 验证方式

```bash
gofmt -w examples/ai/phase3/04_vector_index_and_retrieval/*.go
go test -timeout=60s ./examples/ai/phase3/04_vector_index_and_retrieval
go vet ./examples/ai/phase3/04_vector_index_and_retrieval
```

## 完成标准

- 使用顶部 Qdrant 常量连接阶段 7 Compose 中的真实 Qdrant，完成建集合、写入、查询和删除；内存 Map 不算完成。

- 输入、空结果、维度、Top-K、权限和取消边界得到明确处理。
- 外部依赖错误使用 `%w` 保留原因，不吞掉部分失败。
- 结果可由固定 fixture 重复验证，并记录必要质量或延迟指标。
- 不提交真实 API Key、数据库密码或向量库凭证。

## 暂不实现

- 混合检索、重排和生产级索引切换。
