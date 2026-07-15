# 阶段 6 练习 3：回答事实与引用检查

## 练习目标

检查回答事实、引用和无依据拒答。

## 前置知识

- Golden Dataset、检索指标、回答指标和 Agent 指标需要分层计算。
- 评估必须可重复，优先使用固定 fixture 和确定性 Fake。
- 质量、延迟、Token 和成本需要一起比较。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：从测试集中读取期望事实和允许证据。
2. TODO 2：校验引用 ID 指向实际检索上下文。
3. TODO 3：实现确定性关键词/结构规则检查。
4. TODO 4：对无证据问题验证明确拒答。
5. TODO 5：把 LLM-as-a-Judge 作为可选层，并抽样人工复核。

## 开始练习

```bash
go run ./examples/ai/phase6/03_answer_quality_checks
```

默认骨架不调用付费模型；需要模型结果时先保存离线 fixture。

## 验证方式

```bash
gofmt -w examples/ai/phase6/03_answer_quality_checks/*.go
go test -timeout=60s ./examples/ai/phase6/03_answer_quality_checks
go vet ./examples/ai/phase6/03_answer_quality_checks
```

## 完成标准

- 指标公式由固定小数据集验证，边界 case 有明确期望值。
- 数据缺失、空结果和重复结果不会被静默忽略。
- 报告包含总体、分组和具体失败 case。
- 评估输入、配置和结果带版本，能够本地复现。

## 暂不实现

- 完全自动化的主观质量裁决。
