# 阶段 6 练习 2：检索指标

## 练习目标

自动计算 Retriever 的 Recall@K、MRR 和命中率。

## 前置知识

- Golden Dataset、检索指标、回答指标和 Agent 指标需要分层计算。
- 评估必须可重复，优先使用固定 fixture 和确定性 Fake。
- 质量、延迟、Token 和成本需要一起比较。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：读取 Golden Dataset 与每个查询的排名结果。
2. TODO 2：校验 K、相关文档集合和重复结果。
3. TODO 3：实现 Recall@K、MRR 和 Hit Rate。
4. TODO 4：分别汇总总体与标签分组指标。
5. TODO 5：输出失败 case 和排名证据，避免只给平均数。

## 开始练习

```bash
go run ./examples/ai/phase6/02_retrieval_metrics
```

默认骨架不调用付费模型；需要模型结果时先保存离线 fixture。

## 验证方式

```bash
gofmt -w examples/ai/phase6/02_retrieval_metrics/*.go
go test -timeout=60s ./examples/ai/phase6/02_retrieval_metrics
go vet ./examples/ai/phase6/02_retrieval_metrics
```

## 完成标准

- 指标公式由固定小数据集验证，边界 case 有明确期望值。
- 数据缺失、空结果和重复结果不会被静默忽略。
- 报告包含总体、分组和具体失败 case。
- 评估输入、配置和结果带版本，能够本地复现。

## 暂不实现

- nDCG 和在线点击指标。
