# HTTP API 规范

## Router

总入口是 `router.InitTotalRouter`，使用 `gin.New()` 并按以下顺序注册全局中间件：

```text
CorsMiddleware
-> RequestIdMiddleware
-> GinRecovery
-> GinLogger
-> InjectContainerMiddleware
```

顺序有行为意义：Request ID 必须在恢复和请求日志之前生成，Container 必须在 Handler 调用前注入。调整顺序时要验证 Panic、鉴权失败和正常请求的日志字段。

路由按访问边界分组：

- `/api/public`：无需登录。
- `/api/private`：挂载 `AuthMiddleware`。
- `/api/admin`：依次挂载 `AuthMiddleware` 和 `AdminAuthMiddleware`。

模块路由放在独立 `<feature>_router.go` 中，由 `total_router.go` 统一注册。

## Handler 执行顺序

参考 `handler.UserRegister` 和 `handler.UserLogin`：

1. 使用 `ShouldBindJSON` 或明确的 Path/Query/Form 解析输入。
2. 校验输入；失败立即调用统一错误辅助函数并 `return`。
3. 通过 `getContainerOrFail` 获取已注入依赖。
4. 将 `c.Request.Context()` 传给 Service。
5. 使用 `errors.Is` 区分稳定业务错误。
6. 通过 `handler/result.go` 的辅助函数写入统一响应。

Handler 只处理 HTTP 协议、轻量输入校验和错误到响应的映射。密码哈希、事务、库存、权限数据查询等业务规则放在 Service。

## 统一响应契约

响应结构定义在 `internal/response.Result`：

```json
{"code": 200, "msg": "success", "data": {}, "request_id": "可选"}
```

| 辅助函数 | HTTP 状态 | 业务码 | 用途 |
| --- | ---: | ---: | --- |
| `Success` / `SuccessMsg` | 200 | 200 | 成功 |
| `Fail` | 200 | 400 | 现有业务失败兼容路径 |
| `Error` | 200 | 500 | 现有系统失败兼容路径 |
| `ParamError` | 400 | 400 | 参数错误并中止链路 |
| `AuthError` | 401 | 401 | 未认证并中止链路 |
| `ForbiddenError` | 403 | 403 | 无权限并中止链路 |
| `SystemError` | 500 | 500 | Panic/系统错误，携带 request ID |

`Fail` 和 `Error` 的 HTTP 200 是现有兼容行为，不要在无 API 迁移说明时批量改变。新增协议级错误优先使用对应真实 HTTP 状态的辅助函数。

## Context 与身份

- `RequestIdMiddleware` 将 ID 同时写入 Gin Context 和请求 `context.Context`。
- `AuthMiddleware` 解析 Bearer Token，并把 `int64` 用户 ID 写入 `util.CurrentUserId`。
- Handler 读取用户 ID 时必须检查是否存在和类型是否正确，不能沿用 `userID, _ := c.Get(...)` 的忽略式写法。
- 日志和下游调用使用 `c.Request.Context()`，不能在请求链中换成 `context.Background()`。

## 输入与安全边界

- JSON DTO 使用 `binding` 标签表达必填约束；Path、Query 和 Form 的数字转换必须检查错误与范围。
- 不把数据库错误、堆栈、Token、密码或内部路径返回给客户端。
- `AdminAuthMiddleware` 当前仅检查客户端 `Role: admin` Header，这是已知安全缺口。新管理员能力不能依赖该 Header，应从已验证身份或持久化权限中取角色。
- `CorsMiddleware` 当前会回显请求 Origin 并允许凭证。涉及生产部署或安全改造时必须引入明确白名单，不能另建第二套 CORS 中间件。

## 验证

- Handler 测试使用 `httptest` 驱动 Gin Router，覆盖成功、参数错误、未认证、无权限和 Service 失败。
- 认证相关测试必须覆盖缺失 Header、非 Bearer、无效 JWT 和过期 JWT。
- 中间件变更要验证链路是否 `Abort`、是否重复写响应，以及请求日志中的 HTTP/业务状态码。
