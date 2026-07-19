# 阶段 3 练习 5：HNSW 参数调优

## 练习目标

调整 M、efConstruction、efSearch，记录召回率、内存和 P95 延迟的取舍。

## 前置知识

- Eino 组件的 Context 与错误传播。
- 文档、Chunk、Metadata、向量和 Top-K 的基本含义。
- 默认配置只用于本地骨架验证，真实模型或向量库必须显式配置。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：准备固定语料、查询集和真实目标集合。
2. TODO 2：建立多组 HNSW 索引参数实验。
3. TODO 3：在相同 Top-K 下执行查询并记录命中结果与延迟。
4. TODO 4：计算 Recall@K、P50/P95 和索引大小。
5. TODO 5：输出参数对比表，选择满足目标的最小成本配置。

## 开始练习

先修改 `exercise.go` 顶部的 `qdrantBaseURL`、`qdrantAPIKey` 和 `qdrantCollection`。

```bash
go run ./examples/ai/phase3/05_hnsw_tuning
```

骨架默认返回“练习尚未完成”，不会连接模型、数据库或向量库。

## 验证方式

```bash
gofmt -w examples/ai/phase3/05_hnsw_tuning/*.go
go test -timeout=60s ./examples/ai/phase3/05_hnsw_tuning
go vet ./examples/ai/phase3/05_hnsw_tuning
```

## 完成标准

- 每组 `M`、`efConstruction`、`efSearch` 参数都作用于真实 Qdrant/HNSW 集合并重新执行固定查询；模拟延迟不算完成。

- 输入、空结果、维度、Top-K、权限和取消边界得到明确处理。
- 外部依赖错误使用 `%w` 保留原因，不吞掉部分失败。
- 结果可由固定 fixture 重复验证，并记录必要质量或延迟指标。
- 不提交真实 API Key、数据库密码或向量库凭证。

## 暂不实现

- 十亿级容量规划和分布式向量库运维。
