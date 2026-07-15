# 阶段 2：LLM API 与 Prompt 编号练习

## Goal

完成阶段 2 的 11 个编号练习，迁移并保留已有前 6 个练习，新增 ChatTemplate、Stream、结构化 JSON、调用治理和模型厂商适配骨架。

## Parent Task

- 父任务：`.trellis/tasks/07-15-numbered-ai-exercises/`
- 本子任务只负责阶段 2 的目录、骨架、文档和阶段内验证；跨阶段索引与最终全仓检查由父任务负责。

## Exercise Map

| 顺序 | 目录 | 主题 |
| --- | --- | --- |
| 1 | `01_basic_chat` | 基础模型调用 |
| 2 | `02_eino_generate` | Eino ChatModel Generate |
| 3 | `03_prompt_roles` | System 角色与 User Prompt |
| 4 | `04_in_memory_multi_turn` | 程序内多轮对话 |
| 5 | `05_redis_session_history` | Redis 会话历史 |
| 6 | `06_session_history_limits` | 历史截断、TTL 与用户隔离 |
| 7 | `07_chat_template` | Eino ChatTemplate 与变量 |
| 8 | `08_streaming_chat` | ChatModel Stream、EOF 与取消 |
| 9 | `09_structured_json_output` | 固定 JSON、解析与字段校验 |
| 10 | `10_call_governance` | 超时、重试、限流与调用统计 |
| 11 | `11_model_provider_adapter` | 自定义接口与模型配置切换 |

## Requirements

- 阶段内目录使用两位数字连续编号，并与路线图顺序一致。
- 独立可运行 Go 骨架；默认占位配置必须在创建 Redis 或模型客户端前安全退出。
- 每个练习至少包含中文 README，结构包含目标、前置知识、连续 TODO、运行/验证、完成标准和暂不实现。
- 所有新增 Go 函数和方法写简洁中文职责注释；核心 TODO 按真实执行顺序连续编号。
- 默认占位配置不得连接外部服务、泄露真实密钥或修改外部状态。
- 前 6 个目录只迁移路径和更新引用，不重写内容；第 6 个未提交练习的 TODO 1–11、README 与安全退出行为必须完整保留。

## Acceptance Criteria

- [x] 阶段 2 的 11 个编号目录全部存在，序号连续且没有重复目录。
- [x] 每个 README 的主题、TODO 和完成标准与路线图对应项一致。
- [x] 需要编码的骨架可以编译；未完成路径明确失败且默认不访问外部服务。
- [x] 阶段内 Go 文件通过 `gofmt`、目标包测试与 `go vet`。
- [x] 使用 `go-review` 完成阶段 review；涉及 SQL/Schema/查询时追加 `sql-code-review`。

## Out of Scope

- 提供完整参考答案或生产级实现。
- 修改其他阶段练习或现有电商业务语义。
- 使用真实云资源、真实凭证或执行不可逆外部操作。
