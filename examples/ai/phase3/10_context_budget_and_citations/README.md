# 阶段 3 练习 10：上下文预算与引用

## 练习目标

对候选去重、扩展父块、控制 Token，并返回文件名、页码和证据。

## 前置知识

- Eino 组件的 Context 与错误传播。
- 文档、Chunk、Metadata、向量和 Top-K 的基本含义。
- 默认配置只用于本地骨架验证，真实模型或向量库必须显式配置。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：按文档与位置去重相同或高度重叠 Chunk。
2. TODO 2：按需加载父块，并避免重复扩展相同父块。
3. TODO 3：估算或计算 Token，按预算选择上下文。
4. TODO 4：为每段上下文分配稳定引用 ID，保留文件名和页码。
5. TODO 5：验证回答引用只指向实际发送给模型的证据。

## 开始练习

```bash
go run ./examples/ai/phase3/10_context_budget_and_citations
```

骨架默认返回“练习尚未完成”，不会连接模型、数据库或向量库。

## 验证方式

```bash
gofmt -w examples/ai/phase3/10_context_budget_and_citations/*.go
go test -timeout=60s ./examples/ai/phase3/10_context_budget_and_citations
go vet ./examples/ai/phase3/10_context_budget_and_citations
```

## 完成标准

- 输入、空结果、维度、Top-K、权限和取消边界得到明确处理。
- 外部依赖错误使用 `%w` 保留原因，不吞掉部分失败。
- 结果可由固定 fixture 重复验证，并记录必要质量或延迟指标。
- 不提交真实 API Key、数据库密码或向量库凭证。

## 暂不实现

- 长期记忆和复杂上下文压缩。
