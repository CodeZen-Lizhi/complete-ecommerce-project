# 阶段 5 练习 2：Graph 条件路由

## 练习目标

使用 Graph 将知识、商品和订单问题路由到不同节点。

## 前置知识

- Chain、Graph、ReAct 和普通状态机的职责差异。
- Agent 必须有状态、预算、终止条件和明确副作用边界。
- 固定业务规则优先使用确定性代码。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：定义路由分类和每个节点的输入/输出契约。
2. TODO 2：实现确定性优先的路由条件，低置信度进入澄清节点。
3. TODO 3：使用真实 `compose.NewGraph` 注册知识、商品、订单、澄清和结束节点。
4. TODO 4：连接真实分支边并 `Compile` 为 `compose.Runnable`，验证每条分支可达。
5. TODO 5：覆盖未知分类、节点失败和 Context 取消。

## 开始练习

```bash
go run ./examples/ai/phase5/02_graph_routing
```

骨架默认使用本地状态，不调用真实工具或执行外部副作用。

## 验证方式

```bash
gofmt -w examples/ai/phase5/02_graph_routing/*.go
go test -timeout=60s ./examples/ai/phase5/02_graph_routing
go vet ./examples/ai/phase5/02_graph_routing
```

## 完成标准

- `buildRoutingGraph` 返回真实 Eino Runnable；普通 `switch` 不能替代 Graph 编排。

- 每条路径都有最大步骤、预算和终止条件。
- 状态可序列化、错误可诊断，取消不会被吞掉。
- 高风险副作用在确认前不会执行，恢复过程具备幂等保护。
- Fake 能复现成功、节点失败、循环、取消和恢复场景。

## 暂不实现

- ReAct、工具循环和人工确认。
