# 阶段 3 练习 11：索引生命周期

## 练习目标

使用 MySQL 管理导入状态，实现幂等更新、删除、失败重试和最终一致性补偿。

## 前置知识

- Eino 组件的 Context 与错误传播。
- 文档、Chunk、Metadata、向量和 Top-K 的基本含义。
- 默认配置只用于本地骨架验证，真实模型或向量库必须显式配置。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：定义文档版本、状态、内容哈希和索引 ID 契约。
2. TODO 2：创建导入任务，重复请求按幂等键返回同一任务。
3. TODO 3：完成解析、切块、向量写入后再原子更新可用状态。
4. TODO 4：更新时写新版本并切换，删除时清理向量和元数据。
5. TODO 5：记录失败步骤、重试次数和补偿动作，验证中断恢复。

## 开始练习

```bash
go run ./examples/ai/phase3/11_index_lifecycle
```

骨架默认返回“练习尚未完成”，不会连接模型、数据库或向量库。

## 验证方式

```bash
gofmt -w examples/ai/phase3/11_index_lifecycle/*.go
go test -timeout=60s ./examples/ai/phase3/11_index_lifecycle
go vet ./examples/ai/phase3/11_index_lifecycle
```

## 完成标准

- 从 `MYSQL_DSN` 连接真实 MySQL 保存导入状态，并与真实向量索引执行更新、删除、失败重试和补偿；内存 Store 不算完成。

- 输入、空结果、维度、Top-K、权限和取消边界得到明确处理。
- 外部依赖错误使用 `%w` 保留原因，不吞掉部分失败。
- 结果可由固定 fixture 重复验证，并记录必要质量或延迟指标。
- 不提交真实 API Key、数据库密码或向量库凭证。

## 暂不实现

- 跨区域索引复制和零停机大规模重建。
