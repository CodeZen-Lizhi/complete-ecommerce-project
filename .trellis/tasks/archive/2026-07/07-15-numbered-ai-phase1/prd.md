# 阶段 1：Go 后端基础编号练习

## Goal

为现有电商项目建立 7 个按顺序学习的导航型练习目录，帮助学习者在真实 Handler、Service、Repository、Redis、JWT、HTTP Client、并发和测试边界中完成 Go 后端基础训练。

## Parent Task

- 父任务：`.trellis/tasks/07-15-numbered-ai-exercises/`
- 本子任务只负责阶段 1 的目录、骨架、文档和阶段内验证；跨阶段索引与最终全仓检查由父任务负责。

## Exercise Map

| 顺序 | 目录 | 主题 |
| --- | --- | --- |
| 1 | `01_gin_validated_api` | Gin 参数校验与统一响应 |
| 2 | `02_gorm_crud_pagination_transaction` | GORM CRUD、分页、事务与 Explain |
| 3 | `03_redis_cache_aside` | Redis Cache Aside、空值与失效 |
| 4 | `04_jwt_auth_and_authorization` | JWT 登录、鉴权与后端授权 |
| 5 | `05_http_client_resilience` | HTTP Client 超时与非 2xx 处理 |
| 6 | `06_worker_pool` | Worker Pool、取消与错误收集 |
| 7 | `07_backend_tests` | Service、Handler、数据库与失败路径测试 |

## Requirements

- 阶段内目录使用两位数字连续编号，并与路线图顺序一致。
- 导航型目录：README 指向现有项目代码和建议新增测试位置，不复制生产业务实现。
- 每个练习至少包含中文 README，结构包含目标、前置知识、连续 TODO、运行/验证、完成标准和暂不实现。
- 所有新增 Go 函数和方法写简洁中文职责注释；核心 TODO 按真实执行顺序连续编号。
- 默认占位配置不得连接外部服务、泄露真实密钥或修改外部状态。
- 每个目录主要包含 README；需要学习者新增测试的练习可提供最小 `_test.go` 骨架，但不得复制或分叉现有业务代码。

## Acceptance Criteria

- [x] 阶段 1 的 7 个编号目录全部存在，序号连续且没有重复目录。
- [x] 每个 README 的主题、TODO 和完成标准与路线图对应项一致。
- [x] 需要编码的骨架可以编译；未完成路径明确失败且默认不访问外部服务。
- [x] 阶段内 Go 文件通过 `gofmt`、目标包测试与 `go vet`。
- [x] 使用 `go-review` 完成阶段 review；涉及 SQL/Schema/查询时追加 `sql-code-review`。

## Out of Scope

- 提供完整参考答案或生产级实现。
- 修改其他阶段练习或现有电商业务语义。
- 使用真实云资源、真实凭证或执行不可逆外部操作。
