# 后端开发规范

> 适用于当前 `ecommerce` Go 单体服务。规范依据当前源码、`README.md`、`go.mod` 和根目录 `AGENTS.md`，不是通用 Gin/GORM 模板。

## 架构概览

生产请求的主调用链是：

```text
main.go
  -> container
  -> router
  -> middleware / handler
  -> service
  -> repository
  -> GORM / MySQL
```

横切能力位于 `internal/config`、`internal/logger`、`internal/mysql`、`internal/redis`、`internal/response` 和 `util`。当前模块路径为 `ecommerce`，Go 版本以 `go.mod` 中的 `go 1.25` 为准。

## 开发前检查

- 明确改动属于 HTTP、业务、数据访问、基础设施还是独立学习示例。
- 阅读目标层指南和直接相邻层源码，不跨层猜测接口或数据结构。
- 修改公共接口、认证、事务、DTO 或配置前，搜索所有调用方并确认兼容性。
- 看到 `handler/product_handler.go` 等硬编码结果时，将其视为 Demo 占位，不作为新业务实现范例。
- 不提交 API Key、JWT Secret、数据库密码或其他真实凭证。

## 指南索引

| 指南 | 适用范围 |
| --- | --- |
| [目录与分层](./directory-structure.md) | 包职责、依赖方向、新模块落位 |
| [HTTP API](./http-api-guidelines.md) | Gin 路由、中间件、Handler、统一响应 |
| [Service 与 Repository](./service-repository-guidelines.md) | 业务边界、接口、依赖注入、事务协作 |
| [数据库](./database-guidelines.md) | GORM、模型、查询、事务、Schema 约束 |
| [配置与运行时](./configuration-runtime.md) | Viper、启动顺序、MySQL/Redis、停机 |
| [错误处理](./error-handling.md) | 错误包装、分类、响应映射、Panic 恢复 |
| [日志](./logging-guidelines.md) | `slog`、请求日志、字段与敏感信息 |
| [质量](./quality-guidelines.md) | 格式化、测试、验证、Review 门禁 |

## 交付前质量检查

- `gofmt` 已作用于修改过的 Go 文件，且没有无关格式化。
- 至少运行 `go test ./...`；适用时追加 `go vet ./...` 和受影响流程烟测。
- Handler 没有承载业务规则，Service 没有直接拼 HTTP 响应，Repository 没有处理鉴权或业务文案。
- 请求 `context.Context` 已传到外部调用和数据库边界；事务内使用同一个 `tx` 派生的 Repository。
- 所有错误均被处理或明确传播，没有静默成功、吞错或泄露内部错误给客户端。
- 新代码没有复制 Demo 假数据、弱鉴权、明文密码兼容等已知历史行为。

## 当前已知边界

- 生产 Handler、Service、Repository 仍缺少系统化测试；AI 学习示例已有目录契约测试和少量行为测试，不能把示例测试表述成生产业务覆盖。
- `handler/product_handler.go` 及部分用户接口仍返回假数据或仅回显表单。
- `AdminAuthMiddleware` 当前信任客户端 `Role` Header，这是已知安全缺口，不得扩展为新权限模型。
- Repository 现有接口未全面接收 `context.Context`；新代码应补齐传播，旧接口需单独任务迁移。
- 仓库没有迁移脚本或 `AutoMigrate` 流程，不能臆造数据库 Schema 变更方式。
