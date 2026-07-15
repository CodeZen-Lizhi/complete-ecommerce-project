# 阶段 3 练习 3：切块策略对比

## 练习目标

对比递归长度、结构化、Parent-Child 和表格特殊处理的块边界与召回差异。

## 前置知识

- Eino 组件的 Context 与错误传播。
- 文档、Chunk、Metadata、向量和 Top-K 的基本含义。
- 默认配置只用于本地骨架验证，真实模型或向量库必须显式配置。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：定义统一 Chunk 接口和最大长度、overlap 配置。
2. TODO 2：实现递归/长度切块，保证 overlap 小于块大小。
3. TODO 3：按标题、段落和列表实现结构化切块。
4. TODO 4：生成 Parent-Child ID 关系，并对表格/代码块做独立处理。
5. TODO 5：输出块数量、长度分布、边界示例和父子关系。

## 开始练习

```bash
go run ./examples/ai/phase3/03_chunking_strategies
```

骨架默认返回“练习尚未完成”，不会连接模型、数据库或向量库。

## 验证方式

```bash
gofmt -w examples/ai/phase3/03_chunking_strategies/*.go
go test -timeout=60s ./examples/ai/phase3/03_chunking_strategies
go vet ./examples/ai/phase3/03_chunking_strategies
```

## 完成标准

- 输入、空结果、维度、Top-K、权限和取消边界得到明确处理。
- 外部依赖错误使用 `%w` 保留原因，不吞掉部分失败。
- 结果可由固定 fixture 重复验证，并记录必要质量或延迟指标。
- 不提交真实 API Key、数据库密码或向量库凭证。

## 暂不实现

- Semantic Chunking、Late Chunking 和模型驱动切块。
