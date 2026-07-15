# 阶段 3 练习 12：RAG 检索评估

## 练习目标

建立固定测试集，比较切块、Hybrid、Rerank 的 Recall@K、MRR、延迟和成本。

## 前置知识

- Eino 组件的 Context 与错误传播。
- 文档、Chunk、Metadata、向量和 Top-K 的基本含义。
- 默认配置只用于本地骨架验证，真实模型或向量库必须显式配置。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：定义问题、相关 Chunk ID、期望事实和过滤条件。
2. TODO 2：运行不同检索配置并保存原始排名。
3. TODO 3：实现 Recall@K、MRR 和命中率计算。
4. TODO 4：记录 P50/P95、Embedding/Rerank 调用量和成本。
5. TODO 5：输出可重复报告，并在结果缺失时失败而不是跳过。

## 开始练习

```bash
go run ./examples/ai/phase3/12_rag_evaluation
```

骨架默认返回“练习尚未完成”，不会连接模型、数据库或向量库。

## 验证方式

```bash
gofmt -w examples/ai/phase3/12_rag_evaluation/*.go
go test -timeout=60s ./examples/ai/phase3/12_rag_evaluation
go vet ./examples/ai/phase3/12_rag_evaluation
```

## 完成标准

- 输入、空结果、维度、Top-K、权限和取消边界得到明确处理。
- 外部依赖错误使用 `%w` 保留原因，不吞掉部分失败。
- 结果可由固定 fixture 重复验证，并记录必要质量或延迟指标。
- 不提交真实 API Key、数据库密码或向量库凭证。

## 暂不实现

- LLM-as-a-Judge 和在线用户反馈闭环。
