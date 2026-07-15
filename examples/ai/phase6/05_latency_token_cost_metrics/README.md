# 阶段 6 练习 5：延迟、Token 与成本指标

## 练习目标

记录 P50/P95、Token、失败率和成本。

## 前置知识

- Golden Dataset、检索指标、回答指标和 Agent 指标需要分层计算。
- 评估必须可重复，优先使用固定 fixture 和确定性 Fake。
- 质量、延迟、Token 和成本需要一起比较。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：为每次调用记录开始/结束、首 Token、输入/输出 Token 和错误类型。
2. TODO 2：校验样本数量和缺失字段。
3. TODO 3：实现 P50/P95 分位数计算。
4. TODO 4：按模型、Prompt 版本和场景汇总 Token、成本和失败率。
5. TODO 5：同时展示质量指标，避免只优化成本或延迟。

## 开始练习

```bash
go run ./examples/ai/phase6/05_latency_token_cost_metrics
```

默认骨架不调用付费模型；需要模型结果时先保存离线 fixture。

## 验证方式

```bash
gofmt -w examples/ai/phase6/05_latency_token_cost_metrics/*.go
go test -timeout=60s ./examples/ai/phase6/05_latency_token_cost_metrics
go vet ./examples/ai/phase6/05_latency_token_cost_metrics
```

## 完成标准

- 指标公式由固定小数据集验证，边界 case 有明确期望值。
- 数据缺失、空结果和重复结果不会被静默忽略。
- 报告包含总体、分组和具体失败 case。
- 评估输入、配置和结果带版本，能够本地复现。

## 暂不实现

- 完整 FinOps 平台和供应商账单核对。
