# 阶段 6：Eval 编号练习

## Goal

生成 7 个 Eval 练习骨架，用固定数据集评估 Retriever、回答、Tool、Agent、延迟、Token、成本和回归门禁。

## Parent Task

- 父任务：`.trellis/tasks/07-15-numbered-ai-exercises/`
- 本子任务只负责阶段 6 的目录、骨架、文档和阶段内验证；跨阶段索引与最终全仓检查由父任务负责。

## Exercise Map

| 顺序 | 目录 | 主题 |
| --- | --- | --- |
| 1 | `01_golden_dataset` | 问题、目标 Chunk、事实与工具测试集 |
| 2 | `02_retrieval_metrics` | Recall@K、MRR 与命中率 |
| 3 | `03_answer_quality_checks` | 事实、引用与无依据拒答 |
| 4 | `04_tool_and_agent_metrics` | Tool 参数正确率与 Agent 完成率 |
| 5 | `05_latency_token_cost_metrics` | P50/P95、Token、失败率与成本 |
| 6 | `06_regression_comparison` | 模型、Prompt、切块、Hybrid 与 Rerank 对比 |
| 7 | `07_eval_ci_gate` | 核心评估接入 CI 与质量门禁 |

## Requirements

- 阶段内目录使用两位数字连续编号，并与路线图顺序一致。
- 练习应尽量离线可重复，使用小型 JSON/JSONL fixture 和确定性 Fake，避免评估本身依赖付费模型。
- 每个练习至少包含中文 README，结构包含目标、前置知识、连续 TODO、运行/验证、完成标准和暂不实现。
- 所有新增 Go 函数和方法写简洁中文职责注释；核心 TODO 按真实执行顺序连续编号。
- 默认占位配置不得连接外部服务、泄露真实密钥或修改外部状态。
- 所有指标实现都需要可验证 fixture；LLM-as-a-Judge 只作为可选扩展，并明确人工抽样复核边界。

## Acceptance Criteria

- [x] 阶段 6 的 7 个编号目录全部存在，序号连续且没有重复目录。
- [x] 每个 README 的主题、TODO 和完成标准与路线图对应项一致。
- [x] 需要编码的骨架可以编译；未完成路径明确失败且默认不访问外部服务。
- [x] 阶段内 Go 文件通过 `gofmt`、目标包测试与 `go vet`。
- [x] 使用 `go-review` 完成阶段 review；涉及 SQL/Schema/查询时追加 `sql-code-review`。

## Out of Scope

- 提供完整参考答案或生产级实现。
- 修改其他阶段练习或现有电商业务语义。
- 使用真实云资源、真实凭证或执行不可逆外部操作。
