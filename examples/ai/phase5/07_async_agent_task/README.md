# 阶段 5 练习 7：异步 Agent 任务

## 练习目标

把长任务改为 API 创建任务、Worker 执行、数据库保存状态。

## 前置知识

- Chain、Graph、ReAct 和普通状态机的职责差异。
- Agent 必须有状态、预算、终止条件和明确副作用边界。
- 固定业务规则优先使用确定性代码。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：定义任务状态机和合法状态转换。
2. TODO 2：创建任务时写入幂等键和初始状态。
3. TODO 3：Worker 原子领取任务并定期更新进度。
4. TODO 4：保存 Agent Checkpoint、错误和可展示结果。
5. TODO 5：实现查询/取消接口并处理 Worker 崩溃。

## 开始练习

```bash
go run ./examples/ai/phase5/07_async_agent_task
```

骨架默认使用本地状态，不调用真实工具或执行外部副作用。

## 验证方式

```bash
gofmt -w examples/ai/phase5/07_async_agent_task/*.go
go test -timeout=60s ./examples/ai/phase5/07_async_agent_task
go vet ./examples/ai/phase5/07_async_agent_task
```

## 完成标准

- 每条路径都有最大步骤、预算和终止条件。
- 状态可序列化、错误可诊断，取消不会被吞掉。
- 高风险副作用在确认前不会执行，恢复过程具备幂等保护。
- Fake 能复现成功、节点失败、循环、取消和恢复场景。

## 暂不实现

- 分布式调度平台和跨区域队列。
