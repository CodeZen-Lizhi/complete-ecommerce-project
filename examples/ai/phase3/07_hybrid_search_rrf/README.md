# 阶段 3 练习 7：Dense + BM25 + RRF

## 练习目标

并行执行 Dense Vector 与 BM25 检索，再使用 RRF 融合排名。

## 前置知识

- Eino 组件的 Context 与错误传播。
- 文档、Chunk、Metadata、向量和 Top-K 的基本含义。
- 默认配置只用于本地骨架验证，真实模型或向量库必须显式配置。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：为 Dense 和 BM25 定义相同文档 ID 与权限过滤契约。
2. TODO 2：分别执行两路 Top-K 检索并保留原始排名。
3. TODO 3：实现 RRF 分数 1/(k+rank)，校验 k 为正数。
4. TODO 4：按文档 ID 合并、累加分数、稳定排序和去重。
5. TODO 5：对比单路与融合后的命中率和延迟。

## 开始练习

```bash
go run ./examples/ai/phase3/07_hybrid_search_rrf
```

骨架默认返回“练习尚未完成”，不会连接模型、数据库或向量库。

## 验证方式

```bash
gofmt -w examples/ai/phase3/07_hybrid_search_rrf/*.go
go test -timeout=60s ./examples/ai/phase3/07_hybrid_search_rrf
go vet ./examples/ai/phase3/07_hybrid_search_rrf
```

## 完成标准

- 输入、空结果、维度、Top-K、权限和取消边界得到明确处理。
- 外部依赖错误使用 `%w` 保留原因，不吞掉部分失败。
- 结果可由固定 fixture 重复验证，并记录必要质量或延迟指标。
- 不提交真实 API Key、数据库密码或向量库凭证。

## 暂不实现

- 学习排序模型和复杂查询规划。
