# 日志规范

## 日志实现

项目统一使用标准库 `log/slog`。`internal/logger.InitLogConfig` 根据配置创建单例 Logger：

- `log.level` 控制 debug/info/warn/error。
- `log.encoding: json` 使用 JSON Handler，其他值使用 Text Handler。
- `log.add_source` 控制是否增加源码位置。
- 时间格式由 `ReplaceAttr` 统一为毫秒级本地文本。

业务代码通过 `logger.GetLogger()` 获取实例；启动早期配置失败可使用 `slog` 默认 Logger。

## 结构化字段

使用稳定键值字段，不把多个动态值拼进消息字符串：

```go
log.Info("服务器启动成功", "name", appName, "port", port)
```

请求链字段沿用 `middleware.GinLogger`：

- `request_id`
- `method`、`path`、`route`
- `status_code`、`business_code`
- `latency`、`latency_ms`
- `response_size`、`content_length`
- `client_ip`、`user_agent`、`referer`
- 已认证时的 `user_id`

新链路日志优先使用 `LogAttrs` / `InfoContext` / `WarnContext` / `ErrorContext` 并传递请求 ctx。

## 日志级别

- Debug：开发诊断和低价值明细，例如 `UserRepository.FindByID` 的查询结果摘要。
- Info：服务启动/停机、请求完成、业务流程中的重要正常事件。
- Warn：可恢复或客户端原因的异常，例如鉴权头错误、HTTP 4xx、超过 500ms 的慢请求。
- Error：服务不可用、HTTP 5xx、Panic、资源关闭失败或需要处理的系统故障。

不要把正常 NotFound、用户密码错误等预期业务结果重复记录成 Error，除非它代表系统异常。

## 请求日志与恢复日志

- `RequestIdMiddleware` 必须先于 `GinLogger` 和 `GinRecovery` 注册。
- `GinLogger` 在 `c.Next()` 后统一记录 HTTP 状态、业务码、耗时和响应大小。
- `GinRecovery` 记录 panic 类型、错误链和项目栈，并输出有限的 IDE 可点击位置。
- Handler 只有在需要补充业务上下文时记录；不要重复记录请求中间件已经覆盖的信息。

## 敏感信息

禁止记录：

- `Authorization` Header、JWT、Cookie、API Key、数据库 DSN 和密码。
- 登录/注册明文密码、验证码、私钥和 Secret。
- 完整请求/响应体，除非经过明确字段白名单和脱敏。
- 不必要的邮箱、手机号、地址等个人信息。

`handler.UserRegister` 当前会记录 username、email 和 age。扩展这类日志前先确认隐私必要性；不得加入 password。

## 错误记录

- 使用独立 `"error", err` 字段保留结构，不用 `fmt.Sprintf` 扁平化全部信息。
- 对外部调用记录目标类别、耗时、状态码和重试次数，但不记录密钥和完整敏感 URL。
- 若错误将继续向上传播，由边界层统一记录，内层只负责 `%w` 包装。
- 不以日志代替错误返回或用户响应。

## 验证

- 正常请求、4xx、5xx、慢请求和 Panic 各验证一次日志级别与必需字段。
- 检查请求 ID 在鉴权、Handler、恢复和请求完成日志中一致。
- 搜索新增日志字段，确认没有 token、password、secret、authorization、cookie 或 DSN。
