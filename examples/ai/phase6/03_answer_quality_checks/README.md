# 阶段 6 练习 3：回答事实与引用检查

## 练习目标

检查回答事实、引用和无依据拒答。

## 前置知识

- Golden Dataset、检索指标、回答指标和 Agent 指标需要分层计算。
- 评估必须可重复，Judge 输入使用固定数据集，完成证明来自真实模型调用。
- 质量、延迟、Token 和成本需要一起比较。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：从测试集中读取期望事实和允许证据。
2. TODO 2：校验引用 ID 指向实际检索上下文。
3. TODO 3：实现确定性关键词/结构规则检查。
4. TODO 4：对无证据问题验证明确拒答。
5. TODO 5：使用严格 JSON Schema 调用真实 LLM-as-a-Judge，并与人工标签比较一致率和抽样复核结果。

## 开始练习

先修改 `exercise.go` 顶部的 `baseURL`、`apiKey` 和 `modelName`。

```bash
go run ./examples/ai/phase6/03_answer_quality_checks
```

完成 TODO 后该命令会调用真实模型，运行前从练习顶部常量配置本地服务；不要提交真实 API Key。

## 验证方式

```bash
gofmt -w examples/ai/phase6/03_answer_quality_checks/*.go
go test -timeout=60s ./examples/ai/phase6/03_answer_quality_checks
go vet ./examples/ai/phase6/03_answer_quality_checks
```

## 完成标准

- 使用真实 Eino ChatModel 返回 `judgeResult`，记录原始失败，不用 Fake/Mock 分数作为完成证明。
- Judge 结果必须与人工标签比较，且 Judge 不能作为唯一质量门禁。

- 指标公式由固定小数据集验证，边界 case 有明确期望值。
- 数据缺失、空结果和重复结果不会被静默忽略。
- 报告包含总体、分组和具体失败 case。
- 评估输入、配置和结果带版本，能够本地复现。

## 暂不实现

- 完全自动化的主观质量裁决。
