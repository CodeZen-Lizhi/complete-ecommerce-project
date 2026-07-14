# 代码复用指南

## 原则

先复用项目已有的单一入口，再考虑新增抽象。复用的目标是收敛契约和不变量，不是把几行相似代码全部塞进 `util`。

## 修改前搜索

使用 `rg` 搜索准备新增的概念、字段、错误文案、配置键和调用方式：

```bash
rg -n "SuccessMsg|ParamError|SystemError" handler internal
rg -n "Transaction|WithDB" service repository internal
rg -n "CurrentUserId|ResponseCodeKey|ContainerKey" .
rg -n "New.*Service|New.*Repository" container service repository
```

读取定义和全部调用方后再判断是复用、扩展还是替换。

## 已有复用入口

| 能力 | 唯一/优先入口 | 不要新增 |
| --- | --- | --- |
| API 响应 | Handler 内用 `handler/result.go`，实现位于 `internal/response` | 第二套 Result、局部 JSON 错误格式 |
| 数据库事务 | `internal/mysql.Transaction` + `repo.WithDB(tx)` | Service 内手写 Begin/Commit/Rollback |
| 生产依赖组装 | `container.GetInstance` | Handler 内自行构造 Service/Repository |
| 结构化日志 | `internal/logger.GetLogger` / 注入的 `*slog.Logger` | 新日志框架或包内全局 logger |
| MySQL 连接 | `internal/mysql.InitMySQL` | 业务包重复 `gorm.Open` |
| Redis 连接 | `internal/redis.Init` / `Client` | 每次请求创建 Redis Client |
| JWT | `util.GenerateToken` / `ParseToken` | Handler 内重复解析或签名 |
| 用户/响应 Context Key | `util.CurrentUserId`、`util.ResponseCodeKey` | 相同字符串散落多个包 |
| 请求 ID | `middleware.RequestIdMiddleware` | Handler 各自生成 request ID |
| 雪花 ID | `util.Init` / `GenID` | 每个 Model 重设 Epoch 或 Node |

## 何时扩展已有入口

- 新响应类型仍属于统一 API 契约：扩展 `internal/response`，并在 `handler/result.go` 暴露一致包装。
- 新事务业务：复用 `mysql.Transaction`，为 Repository 增加 `WithDB` 或统一 ctx 契约。
- 新 Service：扩展 Container 的字段和组装，不再建第二个容器。
- 新配置：扩展 `internal/config.Config` 和校验，不在业务包直接读 Viper。
- 新日志字段：扩展现有请求中间件或调用点的结构化 attrs，不建立重复请求 Logger。

扩展公共入口前必须检查全部调用方和兼容性，尤其是响应状态码、Context Key 和接口签名。

## 何时可以新增 helper

只有同时满足以下条件时才抽取：

1. 至少两个真实调用点共享同一语义，而不仅是语法相似。
2. helper 能隐藏稳定不变量，例如错误映射、资源关闭或参数白名单。
3. 名称和归属清晰，放入拥有该能力的包，而不是默认扔进 `util`。
4. 抽取后调用方更容易正确使用，并能通过测试覆盖边界。

例如 Product Handler 中分页解析重复，但它们当前都是 Demo。先完成真实分页契约和输入错误语义，再决定是否抽取；不要围绕占位代码提前设计通用分页框架。

## 不应抽取的内容

- 只出现一次、短且语义明确的代码。
- 跨层混合逻辑，例如同时解析 Gin 参数并执行 GORM 查询的 helper。
- 以 `map[string]any` 隐藏明确 DTO/领域类型的通用函数。
- 为未来猜测需求创建的 BaseRepository、万能 CRUD 或全局 Service Locator。
- 仅为了少写几行而吞掉 error、ctx 或事务边界的包装。

## 重复与第二事实源检查

准备新增常量、错误、响应、配置或校验时，检查：

- 是否已经存在稳定的包级常量或哨兵错误。
- 是否会让同一业务规则同时存在于 Handler 和 Service。
- 是否会让数据库约束、GORM 标签和手写校验互相矛盾。
- 是否会让 Viper、环境变量和业务默认值形成多个优先级来源。
- 是否会让 `internal/response` 与局部 `gin.H` 错误结构分叉。

如果已有入口不够用，优先加深该模块的接口，而不是并排新增另一套实现。

## Review 检查

- 新增 helper 是否有真实重复证据和清楚所有权。
- 是否复用了响应、事务、日志、Container、JWT 和 Context Key。
- 是否因复用而丢失类型、ctx、错误链或安全校验。
- 是否复制了 Demo 假数据、忽略错误或弱权限判断。
- 删除新抽象后，调用方是否反而更清楚；如果是，保留直接实现。
