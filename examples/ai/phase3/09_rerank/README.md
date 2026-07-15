# 阶段 3 练习 9：候选 Rerank

## 练习目标

对 Retriever 候选执行重排，对比重排前后的质量与额外耗时。

## 前置知识

- Eino 组件的 Context 与错误传播。
- 文档、Chunk、Metadata、向量和 Top-K 的基本含义。
- 默认配置只用于本地骨架验证，真实模型或向量库必须显式配置。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：限制候选数量和每个候选文本长度。
2. TODO 2：构造 Query-Document 对并调用可替换 Reranker。
3. TODO 3：校验返回数量、文档 ID 和分数范围。
4. TODO 4：按重排分数稳定排序并处理并列。
5. TODO 5：记录 Recall/MRR 变化、P95 增量和失败降级策略。

## 开始练习

```bash
go run ./examples/ai/phase3/09_rerank
```

骨架默认返回“练习尚未完成”，不会连接模型、数据库或向量库。

## 验证方式

```bash
gofmt -w examples/ai/phase3/09_rerank/*.go
go test -timeout=60s ./examples/ai/phase3/09_rerank
go vet ./examples/ai/phase3/09_rerank
```

## 完成标准

- 输入、空结果、维度、Top-K、权限和取消边界得到明确处理。
- 外部依赖错误使用 `%w` 保留原因，不吞掉部分失败。
- 结果可由固定 fixture 重复验证，并记录必要质量或延迟指标。
- 不提交真实 API Key、数据库密码或向量库凭证。

## 暂不实现

- 端到端训练 Reranker 和在线 A/B 平台。
