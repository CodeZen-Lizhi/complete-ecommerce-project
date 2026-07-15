# 阶段 5 练习 8：任务恢复与幂等

## 练习目标

进程重启后恢复任务，并避免重复执行外部副作用。

## 前置知识

- Chain、Graph、ReAct 和普通状态机的职责差异。
- Agent 必须有状态、预算、终止条件和明确副作用边界。
- 固定业务规则优先使用确定性代码。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：启动时扫描可恢复状态和过期租约。
2. TODO 2：从最近 Checkpoint 重建 AgentState。
3. TODO 3：为每个副作用生成稳定幂等键并持久化结果。
4. TODO 4：恢复时先查询幂等记录，再决定跳过或执行。
5. TODO 5：模拟崩溃点，验证重启后不会重复副作用。

## 开始练习

```bash
go run ./examples/ai/phase5/08_task_recovery_idempotency
```

骨架默认使用本地状态，不调用真实工具或执行外部副作用。

## 验证方式

```bash
gofmt -w examples/ai/phase5/08_task_recovery_idempotency/*.go
go test -timeout=60s ./examples/ai/phase5/08_task_recovery_idempotency
go vet ./examples/ai/phase5/08_task_recovery_idempotency
```

## 完成标准

- 每条路径都有最大步骤、预算和终止条件。
- 状态可序列化、错误可诊断，取消不会被吞掉。
- 高风险副作用在确认前不会执行，恢复过程具备幂等保护。
- Fake 能复现成功、节点失败、循环、取消和恢复场景。

## 暂不实现

- Exactly-once 分布式语义和完整 Saga 平台。
