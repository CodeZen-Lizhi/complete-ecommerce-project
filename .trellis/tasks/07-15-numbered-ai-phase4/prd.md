# 阶段 4：Tool Calling 编号练习

## Goal

生成 7 个 Tool Calling 练习骨架，逐步掌握工具定义、参数校验、身份授权、业务查询、多工具编排、失败治理和安全写操作。

## Parent Task

- 父任务：`.trellis/tasks/07-15-numbered-ai-exercises/`
- 本子任务只负责阶段 4 的目录、骨架、文档和阶段内验证；跨阶段索引与最终全仓检查由父任务负责。

## Exercise Map

| 顺序 | 目录 | 主题 |
| --- | --- | --- |
| 1 | `01_local_readonly_tool` | 本地只读工具与 ToolsNode |
| 2 | `02_tool_argument_validation` | 必填、类型、范围与未知字段校验 |
| 3 | `03_context_identity_authorization` | Context 身份、租户与再次授权 |
| 4 | `04_ecommerce_query_tools` | 商品、库存、订单与物流查询工具 |
| 5 | `05_multi_tool_orchestration` | 多工具注册与连续调用 |
| 6 | `06_tool_resilience` | 超时、下游失败、返回上限与有限重试 |
| 7 | `07_safe_write_tool` | 取消订单/退款模拟、确认、幂等与审计 |

## Requirements

- 阶段内目录使用两位数字连续编号，并与路线图顺序一致。
- 本地小接口与 Fake 实现优先；模型参数一律视为不可信输入，写操作只做受控模拟。
- 每个练习至少包含中文 README，结构包含目标、前置知识、连续 TODO、运行/验证、完成标准和暂不实现。
- 所有新增 Go 函数和方法写简洁中文职责注释；核心 TODO 按真实执行顺序连续编号。
- 默认占位配置不得连接外部服务、泄露真实密钥或修改外部状态。
- 涉及电商领域时通过只读接口或 Fake Repository 表达边界，不直接连接生产数据库；写操作必须默认禁用。

## Acceptance Criteria

- [x] 阶段 4 的 7 个编号目录全部存在，序号连续且没有重复目录。
- [x] 每个 README 的主题、TODO 和完成标准与路线图对应项一致。
- [x] 需要编码的骨架可以编译；未完成路径明确失败且默认不访问外部服务。
- [x] 阶段内 Go 文件通过 `gofmt`、目标包测试与 `go vet`。
- [x] 使用 `go-review` 完成阶段 review；涉及 SQL/Schema/查询时追加 `sql-code-review`。

## Out of Scope

- 提供完整参考答案或生产级实现。
- 修改其他阶段练习或现有电商业务语义。
- 使用真实云资源、真实凭证或执行不可逆外部操作。
