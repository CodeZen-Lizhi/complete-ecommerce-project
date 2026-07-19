# 阶段 4 练习 1：本地只读工具与 ToolsNode

## 练习目标

定义一个无副作用的本地查询工具，并通过 Eino Tool/ToolsNode 执行。

## 前置知识

- 模型只生成 Tool Call，应用负责校验、授权和执行。
- 模型参数是不可信输入，Context 中的身份才是可信边界。
- 写操作需要确认、幂等、审计和明确失败处理。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：定义 Tool Schema、名称、用途和只读边界。
2. TODO 2：解析参数并拒绝空值和未知字段。
3. TODO 3：实现本地只读查询，传播 Context 取消。
4. TODO 4：通过 `compose.NewToolNode` 注册真实 Eino `tool.BaseTool`，并执行一次 Tool Call。
5. TODO 5：检查结果大小和错误，不返回内部敏感数据。

## 开始练习

```bash
go run ./examples/ai/phase4/01_local_readonly_tool
```

骨架默认使用本地接口和安全配置，不连接生产数据库，也不执行真实写操作。

## 验证方式

```bash
gofmt -w examples/ai/phase4/01_local_readonly_tool/*.go
go test -timeout=60s ./examples/ai/phase4/01_local_readonly_tool
go vet ./examples/ai/phase4/01_local_readonly_tool
```

## 完成标准

- 实际创建并调用 `*compose.ToolsNode`；只写自定义 `toolExecutor` 不算完成。

- 参数、身份、租户、权限和结果上限都在工具边界验证。
- Context 取消和依赖错误清楚传播，错误不泄露敏感内部信息。
- 只读调用可测试，写操作默认禁用并具备确认、幂等和审计。
- 不让模型自由决定可信身份或绕过固定业务规则。

## 暂不实现

- 真实数据库和写操作。
