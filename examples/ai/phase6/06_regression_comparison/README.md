# 阶段 6 练习 6：回归对比

## 练习目标

对模型、Prompt、切块、Hybrid 和 Rerank 做可重复回归对比。

## 前置知识

- Golden Dataset、检索指标、回答指标和 Agent 指标需要分层计算。
- 评估必须可重复，优先使用固定 fixture 和确定性 Fake。
- 质量、延迟、Token 和成本需要一起比较。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：定义基线和候选配置的不可变版本标识。
2. TODO 2：使用相同数据集和运行参数执行两组实验。
3. TODO 3：计算质量、延迟、Token 和成本差值。
4. TODO 4：设置统计最小样本和允许波动范围。
5. TODO 5：列出改善、退化和结果不确定的 case。

## 开始练习

```bash
go run ./examples/ai/phase6/06_regression_comparison
```

默认骨架不调用付费模型；需要模型结果时先保存离线 fixture。

## 验证方式

```bash
gofmt -w examples/ai/phase6/06_regression_comparison/*.go
go test -timeout=60s ./examples/ai/phase6/06_regression_comparison
go vet ./examples/ai/phase6/06_regression_comparison
```

## 完成标准

- 指标公式由固定小数据集验证，边界 case 有明确期望值。
- 数据缺失、空结果和重复结果不会被静默忽略。
- 报告包含总体、分组和具体失败 case。
- 评估输入、配置和结果带版本，能够本地复现。

## 暂不实现

- 在线 A/B 分流和显著性平台。
