# 阶段 5 练习 5：循环与重复调用检测

## 练习目标

设置最大步骤，检测重复工具调用和无进展死循环。

## 前置知识

- Chain、Graph、ReAct 和普通状态机的职责差异。
- Agent 必须有状态、预算、终止条件和明确副作用边界。
- 固定业务规则优先使用确定性代码。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：为工具名与规范化参数生成稳定调用指纹。
2. TODO 2：统计相同指纹次数和连续无状态变化次数。
3. TODO 3：在达到最大步骤或重复阈值时终止。
4. TODO 4：返回可诊断错误和最后状态，不伪造成功答案。
5. TODO 5：用确定性 Fake 构造无限循环并验证及时停止。

## 开始练习

```bash
go run ./examples/ai/phase5/05_loop_detection
```

骨架默认使用本地状态，不调用真实工具或执行外部副作用。

## 验证方式

```bash
gofmt -w examples/ai/phase5/05_loop_detection/*.go
go test -timeout=60s ./examples/ai/phase5/05_loop_detection
go vet ./examples/ai/phase5/05_loop_detection
```

## 完成标准

- 每条路径都有最大步骤、预算和终止条件。
- 状态可序列化、错误可诊断，取消不会被吞掉。
- 高风险副作用在确认前不会执行，恢复过程具备幂等保护。
- Fake 能复现成功、节点失败、循环、取消和恢复场景。

## 暂不实现

- 复杂行为学习和动态预算预测。
