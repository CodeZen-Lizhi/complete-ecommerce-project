# 跨层变更指南

## 目标

本项目是分层 Go 单体。跨层修改的主要风险不是“文件漏改”，而是数据、错误、Context、权限和事务在层之间改变语义。先按真实调用链分析，再决定最小修改范围。

## 标准调用链

```text
HTTP Request
-> router/total_router.go（路由组与中间件）
-> middleware（request ID、认证、恢复、日志、Container）
-> handler（输入与响应）
-> service（业务规则与事务）
-> repository（GORM 查询）
-> MySQL
```

返回路径：

```text
Repository error
-> Service 分类/传播
-> Handler 映射
-> internal/response.Result
-> GinLogger 记录 HTTP 状态、业务码和耗时
```

## 新增完整业务能力

以 User 模块为参考，逐项确认：

1. **数据模型**：真实表名、字段、ID、可空性、软删除和敏感字段是什么。
2. **Repository**：查询是否参数化、接收 ctx、处理 NotFound、分页和事务 DB。
3. **Service**：业务不变量、哨兵错误、副作用顺序和事务边界是什么。
4. **Container**：新 Repository/Service 是否在唯一组装点显式注入。
5. **Handler**：DTO、校验、ctx 传播、错误映射和响应 DTO 是否完整。
6. **Router**：路径属于 public/private/admin 哪个边界，中间件顺序是否正确。
7. **日志与测试**：request ID、敏感字段、成功和失败路径是否可验证。

不要从 `handler/product_handler.go` 的假数据直接向下补查询；先建立真实 Model/Repository/Service 契约。

## 修改 API 或 DTO

检查以下消费者：

- `router/*_router.go` 的路径、方法和访问边界。
- Handler 输入结构和 `binding` 标签。
- Service 方法签名及实现。
- Repository 查询参数。
- `internal/response.Result` 与 Handler 响应辅助函数。
- 客户端是否依赖现有 HTTP 200 + 业务失败码兼容行为。
- 日志是否依赖 `business_code` 或用户 ID。

持久化 Model 不应直接作为包含敏感字段的公开响应。`model.User` 包含 Password；新增用户响应必须使用不含 Password 的 DTO。

## 修改认证或权限

认证链路是：

```text
Authorization Header
-> parseBearerToken
-> util.ParseToken
-> util.CurrentUserId
-> private/admin Handler
```

变更 JWT 时同步检查：

- `util.UserClaims`、issuer、subject、audience、有效期和算法白名单。
- `internal/config.JwtConfig`、YAML、校验和本地 Secret 注入。
- `AuthMiddleware` 的错误响应和日志脱敏。
- 路由组是否正确挂载中间件。
- 旧 Token 是否需要兼容；当前 `ParseToken` 对缺少部分标准字段的旧 Token 有兼容逻辑。

管理员权限当前依赖客户端 `Role` Header，是已知缺口。任何新增管理员写操作都必须先建立可信角色来源，不能扩大该模式。

## 修改事务或数据写入

检查数据是否必须原子完成。如果需要事务：

```text
Handler ctx
-> Service
-> mysql.Transaction(ctx)
-> repo.WithDB(tx)
-> 所有读写均使用事务 repo
```

确认：

- 任一步失败都会返回非 nil error。
- 没有混用原 Repository 和事务 Repository。
- 外部网络调用是否应该放在事务外，避免长时间持锁。
- 重试是否可能造成重复创建或重复扣减。
- `Updates(struct)` 是否错误忽略零值。

## 修改配置或运行时资源

配置字段至少跨越：

```text
configs/config.<env>.yaml
<-> internal/config.Config/mapstructure
-> validateConfig
-> logger/mysql/redis/main/util 消费者
```

不能因为 Viper WatchConfig 存在，就假设已初始化的 `sync.Once` Logger、Redis 或 MySQL 自动刷新。需要动态生效时，把资源切换和并发可见性作为独立设计。

新增进程资源还要检查 `main.go` 的启动失败、Close、退出信号和 Shutdown 超时。

## 修改统一响应或日志

`internal/response.writeResult` 同时影响：

- 客户端 JSON 契约。
- HTTP 状态码。
- Gin 链路是否 Abort。
- `util.ResponseCodeKey`。
- `GinLogger` 的 `business_code`。
- `GinRecovery` 的系统错误响应。

修改它之前必须读取所有辅助函数和中间件调用方，不能只验证一个 Handler。

## 交付前问题

- 输入从哪里进入，在哪一层完成可信校验？
- ctx 是否从 HTTP 一直传到 DB/Redis/HTTP/模型调用？
- 哪一层拥有业务不变量，哪一层只做协议适配？
- 错误能否被上层 `errors.Is/As` 分类，客户端会看到什么？
- 事务、幂等、软删除和权限是否在所有分支一致？
- 日志是否包含 request ID 且不泄露敏感数据？
- 测试是否覆盖成功、拒绝、依赖失败、取消和边界值？
