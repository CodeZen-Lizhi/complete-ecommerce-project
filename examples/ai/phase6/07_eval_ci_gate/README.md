# 阶段 6 练习 7：CI 质量门禁

## 练习目标

将核心评估接入 CI，质量明显下降时阻止发布。

## 前置知识

- Golden Dataset、检索指标、回答指标和 Agent 指标需要分层计算。
- 评估必须可重复，优先使用固定 fixture 和确定性 Fake。
- 质量、延迟、Token 和成本需要一起比较。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：选择离线、快速、确定性的核心 case。
2. TODO 2：定义绝对阈值和相对基线退化阈值。
3. TODO 3：区分代码失败、数据缺失和质量退化。
4. TODO 4：生成机器可读报告与人类摘要。
5. TODO 5：保证失败返回非零退出码，并提供本地复现命令。

## 开始练习

```bash
go run ./examples/ai/phase6/07_eval_ci_gate
```

默认骨架不调用付费模型；需要模型结果时先保存离线 fixture。

## 验证方式

```bash
gofmt -w examples/ai/phase6/07_eval_ci_gate/*.go
go test -timeout=60s ./examples/ai/phase6/07_eval_ci_gate
go vet ./examples/ai/phase6/07_eval_ci_gate
```

## 完成标准

- 指标公式由固定小数据集验证，边界 case 有明确期望值。
- 数据缺失、空结果和重复结果不会被静默忽略。
- 报告包含总体、分组和具体失败 case。
- 评估输入、配置和结果带版本，能够本地复现。

## 暂不实现

- 在 CI 中调用付费模型和大规模全量评估。
