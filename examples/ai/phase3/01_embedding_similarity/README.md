# 阶段 3 练习 1：Embedding、余弦相似度与 Top-K

## 练习目标

调用 Embedding 接口得到向量，手工实现余弦相似度并返回最相关的 Top-K 文本。

## 前置知识

- Eino 组件的 Context 与错误传播。
- 文档、Chunk、Metadata、向量和 Top-K 的基本含义。
- 配置使用 `OPENAI_BASE_URL`、`EMBEDDING_MODEL` 和 `OPENAI_API_KEY`；源码不保存真实密钥。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：校验文本集合、查询文本、向量维度和 TopK。
2. TODO 2：创建 Embedder，并分别生成文档向量与查询向量。
3. TODO 3：实现余弦相似度，处理零向量和维度不一致。
4. TODO 4：为每个文档计算分数，稳定排序并截取 Top-K。
5. TODO 5：打印文本、分数、维度和耗时，不输出 API Key。

## 开始练习

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
