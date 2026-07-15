# 阶段 6 练习 1：Golden Dataset

## 练习目标

建立包含问题、目标 Chunk、期望事实和期望工具的固定测试集。

## 前置知识

- Golden Dataset、检索指标、回答指标和 Agent 指标需要分层计算。
- 评估必须可重复，优先使用固定 fixture 和确定性 Fake。
- 质量、延迟、Token 和成本需要一起比较。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：定义版本化 JSONL Schema 和唯一 case ID。
2. TODO 2：加入查询、相关 Chunk ID、期望事实、允许/禁止工具和标签。
3. TODO 3：加载时严格校验未知字段、重复 ID 和空目标。
4. TODO 4：按主题、难度、租户和风险分层抽样。
5. TODO 5：记录数据集版本和变更原因，保证评估可复现。

## 开始练习

```bash
go run ./examples/ai/phase6/01_golden_dataset
```

默认骨架不调用付费模型；需要模型结果时先保存离线 fixture。

## 验证方式

```bash
gofmt -w examples/ai/phase6/01_golden_dataset/*.go
go test -timeout=60s ./examples/ai/phase6/01_golden_dataset
go vet ./examples/ai/phase6/01_golden_dataset
```

## 完成标准

- 指标公式由固定小数据集验证，边界 case 有明确期望值。
- 数据缺失、空结果和重复结果不会被静默忽略。
- 报告包含总体、分组和具体失败 case。
- 评估输入、配置和结果带版本，能够本地复现。

## 暂不实现

- 自动生成大规模测试集和纯 LLM 评审。
