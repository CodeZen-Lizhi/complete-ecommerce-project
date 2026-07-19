# 阶段 3 练习 8：问题改写、Multi-Query 与 Metadata Filter

## 练习目标

把多轮追问改写为独立查询，生成多个检索查询并强制应用 Metadata Filter。

## 前置知识

- Eino 组件的 Context 与错误传播。
- 文档、Chunk、Metadata、向量和 Top-K 的基本含义。
- 默认配置只用于本地骨架验证，真实模型或向量库必须显式配置。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：从受控历史中提取当前问题需要的指代信息。
2. TODO 2：生成独立查询，并拒绝改变用户原始意图。
3. TODO 3：生成有限数量 Multi-Query，去重并限制并发。
4. TODO 4：在每一路检索强制加入租户、权限、时间和文档类型过滤。
5. TODO 5：融合结果并记录每个查询的命中贡献。

## 开始练习

先修改 `exercise.go` 顶部的模型与 Qdrant 占位常量。

```bash
go run ./examples/ai/phase3/08_query_rewrite_and_filter
```

骨架默认返回“练习尚未完成”，不会连接模型、数据库或向量库。

## 验证方式

```bash
gofmt -w examples/ai/phase3/08_query_rewrite_and_filter/*.go
go test -timeout=60s ./examples/ai/phase3/08_query_rewrite_and_filter
go vet ./examples/ai/phase3/08_query_rewrite_and_filter
```

## 完成标准

- 输入、空结果、维度、Top-K、权限和取消边界得到明确处理。
- 外部依赖错误使用 `%w` 保留原因，不吞掉部分失败。
- 结果可由固定 fixture 重复验证，并记录必要质量或延迟指标。
- 不提交真实 API Key、数据库密码或向量库凭证。

## 暂不实现

- HyDE、自动查询路由和跨语言检索。
