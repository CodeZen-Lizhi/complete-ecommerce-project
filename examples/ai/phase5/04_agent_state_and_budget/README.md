# 阶段 5 练习 4：Agent 状态与预算

## 练习目标

保存当前步骤、工具结果、错误、Token 和预算状态。

## 前置知识

- Chain、Graph、ReAct 和普通状态机的职责差异。
- Agent 必须有状态、预算、终止条件和明确副作用边界。
- 固定业务规则优先使用确定性代码。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：定义可序列化 AgentState 和状态版本。
2. TODO 2：每一步执行前检查最大步骤、Token 和成本预算。
3. TODO 3：记录工具输入摘要、结果引用和错误分类。
4. TODO 4：状态更新采用复制或受控写入，避免隐式共享。
5. TODO 5：在成功、失败和取消时保存最终状态。

## 开始练习

```bash
go run ./examples/ai/phase5/04_agent_state_and_budget
```

骨架默认使用本地状态，不调用真实工具或执行外部副作用。

## 验证方式

```bash
gofmt -w examples/ai/phase5/04_agent_state_and_budget/*.go
go test -timeout=60s ./examples/ai/phase5/04_agent_state_and_budget
go vet ./examples/ai/phase5/04_agent_state_and_budget
```

## 完成标准

- 每条路径都有最大步骤、预算和终止条件。
- 状态可序列化、错误可诊断，取消不会被吞掉。
- 高风险副作用在确认前不会执行，恢复过程具备幂等保护。
- Fake 能复现成功、节点失败、循环、取消和恢复场景。

## 暂不实现

- 跨区域状态复制和精确计费平台。
