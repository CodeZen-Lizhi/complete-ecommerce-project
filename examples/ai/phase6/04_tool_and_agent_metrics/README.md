# 阶段 6 练习 4：Tool 与 Agent 指标

## 练习目标

统计 Tool 参数正确率和 Agent 任务完成率。

## 前置知识

- Golden Dataset、检索指标、回答指标和 Agent 指标需要分层计算。
- 评估必须可重复，优先使用固定 fixture 和确定性 Fake。
- 质量、延迟、Token 和成本需要一起比较。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：定义期望工具、参数约束和任务终态。
2. TODO 2：对比实际 Tool Call 名称和规范化参数。
3. TODO 3：统计选择率、参数正确率、非法调用率和完成率。
4. TODO 4：记录步骤数、重复调用和终止原因。
5. TODO 5：输出失败轨迹，区分模型、工具和环境问题。

## 开始练习

```bash
go run ./examples/ai/phase6/04_tool_and_agent_metrics
```

默认骨架不调用付费模型；需要模型结果时先保存离线 fixture。

## 验证方式

```bash
gofmt -w examples/ai/phase6/04_tool_and_agent_metrics/*.go
go test -timeout=60s ./examples/ai/phase6/04_tool_and_agent_metrics
go vet ./examples/ai/phase6/04_tool_and_agent_metrics
```

## 完成标准

- 指标公式由固定小数据集验证，边界 case 有明确期望值。
- 数据缺失、空结果和重复结果不会被静默忽略。
- 报告包含总体、分组和具体失败 case。
- 评估输入、配置和结果带版本，能够本地复现。

## 暂不实现

- 线上用户满意度归因。
