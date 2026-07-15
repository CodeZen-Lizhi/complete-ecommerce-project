# 阶段 5 练习 6：Interrupt、Checkpoint 与 Resume

## 练习目标

在高风险节点中断，保存 Checkpoint，确认后恢复执行。

## 前置知识

- Chain、Graph、ReAct 和普通状态机的职责差异。
- Agent 必须有状态、预算、终止条件和明确副作用边界。
- 固定业务规则优先使用确定性代码。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：定义高风险节点和中断原因。
2. TODO 2：中断前保存版本化 Checkpoint 与待确认动作摘要。
3. TODO 3：返回确认标识，不在中断阶段执行副作用。
4. TODO 4：恢复时校验确认人、状态版本和过期时间。
5. TODO 5：保证 Resume 幂等，并覆盖拒绝、过期和重复确认。

## 开始练习

```bash
go run ./examples/ai/phase5/06_interrupt_checkpoint_resume
```

骨架默认使用本地状态，不调用真实工具或执行外部副作用。

## 验证方式

```bash
gofmt -w examples/ai/phase5/06_interrupt_checkpoint_resume/*.go
go test -timeout=60s ./examples/ai/phase5/06_interrupt_checkpoint_resume
go vet ./examples/ai/phase5/06_interrupt_checkpoint_resume
```

## 完成标准

- 每条路径都有最大步骤、预算和终止条件。
- 状态可序列化、错误可诊断，取消不会被吞掉。
- 高风险副作用在确认前不会执行，恢复过程具备幂等保护。
- Fake 能复现成功、节点失败、循环、取消和恢复场景。

## 暂不实现

- 真实支付审批和多级工作流平台。
