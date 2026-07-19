# 阶段 5 练习 3：ReAct Agent

## 练习目标

让 ChatModel 选择工具，ToolsNode 执行并把结果返回模型。

## 前置知识

- Chain、Graph、ReAct 和普通状态机的职责差异。
- Agent 必须有状态、预算、终止条件和明确副作用边界。
- 固定业务规则优先使用确定性代码。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：准备只读工具及严格 Schema。
2. TODO 2：以不可变方式为真实 `model.ToolCallingChatModel` 配置 `compose.ToolsNodeConfig`。
3. TODO 3：通过 `react.NewAgent` 创建真实 Eino ReAct Agent。
4. TODO 4：关联 Tool Call 与 Tool Result，并把结果送回模型。
5. TODO 5：设置最大步骤并验证不调用、单次调用和多次调用。

## 开始练习

```bash
go run ./examples/ai/phase5/03_react_agent
```

骨架默认使用本地状态，不调用真实工具或执行外部副作用。

## 验证方式

```bash
gofmt -w examples/ai/phase5/03_react_agent/*.go
go test -timeout=60s ./examples/ai/phase5/03_react_agent
go vet ./examples/ai/phase5/03_react_agent
```

## 完成标准

- 实际运行 `*react.Agent`；只实现本地 `reactModel/reactToolsNode` 循环不算完成。

- 每条路径都有最大步骤、预算和终止条件。
- 状态可序列化、错误可诊断，取消不会被吞掉。
- 高风险副作用在确认前不会执行，恢复过程具备幂等保护。
- Fake 能复现成功、节点失败、循环、取消和恢复场景。

## 暂不实现

- 写工具、Checkpoint 和长期记忆。
