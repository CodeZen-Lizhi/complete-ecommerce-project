# 阶段 4 练习 5：多工具注册与连续调用

## 练习目标

注册多个工具，处理不调用、单次调用和连续调用。

## 前置知识

- 模型只生成 Tool Call，应用负责校验、授权和执行。
- 模型参数是不可信输入，Context 中的身份才是可信边界。
- 写操作需要确认、幂等、审计和明确失败处理。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：为工具提供互不混淆的名称和描述。
2. TODO 2：通过真实 `ToolCallingChatModel.WithTools` 和 `compose.NewToolNode` 不可变绑定工具。
3. TODO 3：真实调用模型，处理不调用工具、调用一个工具和连续调用多个工具，并把结果返回模型。
4. TODO 4：把每次 Tool Result 与对应 Tool Call ID 关联。
5. TODO 5：设置最大工具调用次数并检测重复调用。

## 开始练习

```bash
go run ./examples/ai/phase4/05_multi_tool_orchestration
```

完成 TODO 后会产生真实模型调用和只读工具调用；不得连接生产数据库，也不得执行真实写操作。

## 验证方式

```bash
gofmt -w examples/ai/phase4/05_multi_tool_orchestration/*.go
go test -timeout=60s ./examples/ai/phase4/05_multi_tool_orchestration
go vet ./examples/ai/phase4/05_multi_tool_orchestration
```

## 完成标准

- 实际执行 ToolCallingChatModel → ToolsNode → Tool Result → ToolCallingChatModel 闭环；本地 `modelTurn` 模拟不能作为完成证明。

- 参数、身份、租户、权限和结果上限都在工具边界验证。
- Context 取消和依赖错误清楚传播，错误不泄露敏感内部信息。
- 只读调用可测试，写操作默认禁用并具备确认、幂等和审计。
- 不让模型自由决定可信身份或绕过固定业务规则。

## 暂不实现

- 完整 ReAct Agent 和长期状态。
