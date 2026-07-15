# 阶段 5 练习 1：固定顺序 Chain

## 练习目标

使用 Eino Chain 完成输入规范化、Prompt、模型调用和结果解析的固定流程。

## 前置知识

- Chain、Graph、ReAct 和普通状态机的职责差异。
- Agent 必须有状态、预算、终止条件和明确副作用边界。
- 固定业务规则优先使用确定性代码。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：定义每个节点的输入/输出类型和唯一职责。
2. TODO 2：按确定顺序组装 Chain，不为固定规则引入模型决策。
3. TODO 3：传播 Context、节点错误和取消。
4. TODO 4：编译 Chain 后使用固定输入执行。
5. TODO 5：记录每个节点耗时并验证空输入和中间失败。

## 开始练习

```bash
go run ./examples/ai/phase5/01_chain_workflow
```

骨架默认使用本地状态，不调用真实工具或执行外部副作用。

## 验证方式

```bash
gofmt -w examples/ai/phase5/01_chain_workflow/*.go
go test -timeout=60s ./examples/ai/phase5/01_chain_workflow
go vet ./examples/ai/phase5/01_chain_workflow
```

## 完成标准

- 每条路径都有最大步骤、预算和终止条件。
- 状态可序列化、错误可诊断，取消不会被吞掉。
- 高风险副作用在确认前不会执行，恢复过程具备幂等保护。
- Fake 能复现成功、节点失败、循环、取消和恢复场景。

## 暂不实现

- 条件分支、循环和持久化状态。
