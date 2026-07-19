# 阶段 3 练习 6：Parent-Child Retrieval

## 练习目标

使用小子块精确召回，命中后返回更完整的父块上下文。

## 前置知识

- Eino 组件的 Context 与错误传播。
- 文档、Chunk、Metadata、向量和 Top-K 的基本含义。
- 默认配置只用于本地骨架验证，真实模型或向量库必须显式配置。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：导入时生成稳定 parent ID 和 child ID。
2. TODO 2：只为 child 建立向量索引，父块保存于可查询存储。
3. TODO 3：检索 child 后按 parent ID 批量加载父块。
4. TODO 4：按最佳子块分数合并父块并去重。
5. TODO 5：验证父块权限继承、缺失父块和一父多子顺序。

## 开始练习

先修改 `exercise.go` 顶部的 `qdrantBaseURL`、`qdrantAPIKey` 和 `qdrantCollection`。

```bash
go run ./examples/ai/phase3/06_parent_child_retrieval
```

骨架默认返回“练习尚未完成”，不会连接模型、数据库或向量库。

## 验证方式

```bash
gofmt -w examples/ai/phase3/06_parent_child_retrieval/*.go
go test -timeout=60s ./examples/ai/phase3/06_parent_child_retrieval
go vet ./examples/ai/phase3/06_parent_child_retrieval
```

## 完成标准

- 输入、空结果、维度、Top-K、权限和取消边界得到明确处理。
- 外部依赖错误使用 `%w` 保留原因，不吞掉部分失败。
- 结果可由固定 fixture 重复验证，并记录必要质量或延迟指标。
- 不提交真实 API Key、数据库密码或向量库凭证。

## 暂不实现

- 多级父子树和跨文档上下文拼接。
