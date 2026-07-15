# 阶段 7：安全监控与部署编号练习

## Goal

生成 8 个安全、可观测性和部署练习骨架，覆盖 Secret、Prompt Injection、限流、Trace、Metrics、Compose、灰度回滚和故障降级。

## Parent Task

- 父任务：`.trellis/tasks/07-15-numbered-ai-exercises/`
- 本子任务只负责阶段 7 的目录、骨架、文档和阶段内验证；跨阶段索引与最终全仓检查由父任务负责。

## Exercise Map

| 顺序 | 目录 | 主题 |
| --- | --- | --- |
| 1 | `01_secret_and_log_redaction` | Secret 管理与日志脱敏 |
| 2 | `02_prompt_injection_defense` | Prompt Injection 与权限边界 |
| 3 | `03_rate_limit_and_concurrency` | 接口、模型与 Embedding 限流 |
| 4 | `04_end_to_end_tracing` | 模型、Retriever、Rerank、Tool 与 Agent Trace |
| 5 | `05_metrics_and_alerts` | 耗时、Token、成本、错误率与队列告警 |
| 6 | `06_docker_compose_stack` | 应用、MySQL、Redis 与向量库 Compose |
| 7 | `07_health_gray_release_rollback` | 健康检查、灰度与回滚 |
| 8 | `08_dependency_failure_fallback` | 模型或向量库故障与降级 |

## Requirements

- 阶段内目录使用两位数字连续编号，并与路线图顺序一致。
- 默认离线或本地安全模式；任何凭证、真实发布、外部告警和破坏性故障注入都必须由学习者显式启用。
- 每个练习至少包含中文 README，结构包含目标、前置知识、连续 TODO、运行/验证、完成标准和暂不实现。
- 所有新增 Go 函数和方法写简洁中文职责注释；核心 TODO 按真实执行顺序连续编号。
- 默认占位配置不得连接外部服务、泄露真实密钥或修改外部状态。
- 部署练习生成安全模板和验证清单，不执行真实发布；Compose 配置不得包含真实密钥，故障注入只针对本地环境。

## Acceptance Criteria

- [x] 阶段 7 的 8 个编号目录全部存在，序号连续且没有重复目录。
- [x] 每个 README 的主题、TODO 和完成标准与路线图对应项一致。
- [x] 需要编码的骨架可以编译；未完成路径明确失败且默认不访问外部服务。
- [x] 阶段内 Go 文件通过 `gofmt`、目标包测试与 `go vet`。
- [x] 使用 `go-review` 完成阶段 review；涉及 SQL/Schema/查询时追加 `sql-code-review`。

## Out of Scope

- 提供完整参考答案或生产级实现。
- 修改其他阶段练习或现有电商业务语义。
- 使用真实云资源、真实凭证或执行不可逆外部操作。
