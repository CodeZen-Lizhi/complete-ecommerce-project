# 阶段 5：Agent 工作流编号练习

## Goal

生成 8 个 Agent 工作流练习骨架，覆盖 Chain、Graph、ReAct、状态与预算、循环控制、Interrupt/Checkpoint/Resume、异步任务和幂等恢复。

## Parent Task

- 父任务：`.trellis/tasks/07-15-numbered-ai-exercises/`
- 本子任务只负责阶段 5 的目录、骨架、文档和阶段内验证；跨阶段索引与最终全仓检查由父任务负责。

## Exercise Map

| 顺序 | 目录 | 主题 |
| --- | --- | --- |
| 1 | `01_chain_workflow` | 固定顺序 Chain |
| 2 | `02_graph_routing` | 知识、商品与订单 Graph 路由 |
| 3 | `03_react_agent` | ChatModel、ToolsNode 与 ReAct |
| 4 | `04_agent_state_and_budget` | 步骤、结果、错误、Token 与预算状态 |
| 5 | `05_loop_detection` | 最大步骤、重复调用与死循环检测 |
| 6 | `06_interrupt_checkpoint_resume` | Interrupt、Checkpoint 与 Resume |
| 7 | `07_async_agent_task` | API 创建任务、Worker 执行与状态持久化 |
| 8 | `08_task_recovery_idempotency` | 重启恢复与幂等副作用 |

## Requirements

- 阶段内目录使用两位数字连续编号，并与路线图顺序一致。
- 优先确定性流程和内存 Fake；高风险节点必须显式中断，默认不执行真实副作用。
- 每个练习至少包含中文 README，结构包含目标、前置知识、连续 TODO、运行/验证、完成标准和暂不实现。
- 所有新增 Go 函数和方法写简洁中文职责注释；核心 TODO 按真实执行顺序连续编号。
- 默认占位配置不得连接外部服务、泄露真实密钥或修改外部状态。
- Checkpoint 与任务持久化先定义可替换接口和内存实现；需要 Redis/MySQL 时只在 README 中给出后续扩展点。

## Acceptance Criteria

- [x] 阶段 5 的 8 个编号目录全部存在，序号连续且没有重复目录。
- [x] 每个 README 的主题、TODO 和完成标准与路线图对应项一致。
- [x] 需要编码的骨架可以编译；未完成路径明确失败且默认不访问外部服务。
- [x] 阶段内 Go 文件通过 `gofmt`、目标包测试与 `go vet`。
- [x] 使用 `go-review` 完成阶段 review；涉及 SQL/Schema/查询时追加 `sql-code-review`。

## Out of Scope

- 提供完整参考答案或生产级实现。
- 修改其他阶段练习或现有电商业务语义。
- 使用真实云资源、真实凭证或执行不可逆外部操作。
