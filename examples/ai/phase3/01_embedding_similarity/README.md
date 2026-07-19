# 阶段 3 练习 1：Embedding、余弦相似度与 Top-K

## 练习目标

调用 Embedding 接口得到向量，手工实现余弦相似度并返回最相关的 Top-K 文本。

## 前置知识

- Eino 组件的 Context 与错误传播。
- 文档、Chunk、Metadata、向量和 Top-K 的基本含义。
- 模型连接通过 `exercise.go` 顶部的 `baseURL`、`apiKey` 和 `modelName` 占位常量配置。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：实现 `parseEmbeddingInput`，校验查询、候选文本和 Top-K。
2. TODO 2：实现 `loadEmbeddingConfig`，读取并校验模型配置。
3. TODO 3：实现 `embedInputs`，分别生成候选向量和查询向量。
4. TODO 4：实现 `cosineSimilarity`，处理零向量和维度不一致。
5. TODO 5：实现 `rankTopK`，稳定排序并截取 Top-K。
6. TODO 6：实现 `reportEmbeddingResults`，打印分数、维度和耗时且不输出 API Key。

## 开始练习

先修改 `exercise.go` 顶部的模型占位常量；本地练习可以填写真实值，但不要提交真实 API Key。

```bash
go run ./examples/ai/phase3/01_embedding_similarity \
  --text "Go 的 context 可以传播取消信号" \
  --candidate "Context 用于取消和超时" \
  --candidate "Goroutine 是轻量线程" \
  --top-k 1
```

骨架默认返回“练习尚未完成”，不会连接模型、数据库或向量库。

## 验证方式

```bash
gofmt -w examples/ai/phase3/01_embedding_similarity/*.go
go test -timeout=60s ./examples/ai/phase3/01_embedding_similarity
go vet ./examples/ai/phase3/01_embedding_similarity
```

## 完成标准

- 输入、空结果、维度、Top-K、权限和取消边界得到明确处理。
- 外部依赖错误使用 `%w` 保留原因，不吞掉部分失败。
- 结果可由固定 fixture 重复验证，并记录必要质量或延迟指标。
- 不提交真实 API Key、数据库密码或向量库凭证。

## 暂不实现

- 向量数据库、HNSW 和 Rerank。
