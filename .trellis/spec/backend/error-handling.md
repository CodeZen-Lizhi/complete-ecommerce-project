# 错误处理规范

## 基本规则

- Go 错误显式返回和检查，不使用异常式控制流。
- 增加操作上下文时使用 `fmt.Errorf("...: %w", err)`，保留错误链。
- 判断包装错误使用 `errors.Is` / `errors.As`，不匹配错误字符串。
- 函数应返回失败，不用空对象、零值或成功响应掩盖错误。
- 同一错误只在合适的边界记录一次，避免每层重复 log 后再 return。

## 各层职责

### Repository

`repository.wrapRepoError` 为数据库错误补充操作名和源码位置，并使用 `%w` 保留 GORM 错误。Repository 不生成用户文案，不决定 HTTP 状态。

`IsExist` 将 `gorm.ErrRecordNotFound` 映射为 `(false, nil)`，因为“不存在”是该方法的正常布尔结果；`FindByUsername` 则保留 NotFound 错误。新增方法必须在接口语义中明确采用哪一种。

### Service

Service 定义稳定业务错误，例如 `service.ErrUserAlreadyExists`，供 Handler 使用 `errors.Is` 分类。底层故障保留错误链返回，不泄露到客户端。

`UserService.Login` 把“用户不存在”和“密码不匹配”统一为相同文案，避免泄露用户名是否存在。

### Handler

Handler 将参数、业务、鉴权和系统错误映射为 `handler/result.go` 中的统一响应辅助函数，然后立即 `return`。客户端只看到稳定文案，不看到 SQL、堆栈、绝对路径或内部错误链。

### 进程入口

`main.go` 在配置或基础设施初始化失败时记录错误并终止启动流程；资源初始化成功后注册关闭逻辑。不要让服务在 DB/Redis 初始化失败后继续对外宣称可用。

## Panic 恢复

`middleware.GinRecovery` 是 HTTP 边界的最后防线：

- 将 panic 值转换为 error。
- 记录错误链、request ID、方法、路径和过滤后的业务栈。
- 若响应尚未写入，调用 `response.SystemError` 返回统一 500。
- 若 Handler 已写响应，不重复写入。

Panic Recovery 不能替代正常错误处理。预期失败必须返回 error；不要主动 panic 来处理参数、NotFound 或外部服务错误。

## Context 错误

- `context.Canceled` 表示上游取消，通常不应重试。
- `context.DeadlineExceeded` 表示超时；是否重试取决于操作幂等性和整体 deadline。
- 请求链中的 DB、Redis、HTTP 和模型调用必须使用上游 ctx，确保取消能够传播。
- 错误包装后仍使用 `errors.Is(err, context.DeadlineExceeded)` 判断。

## 当前兼容行为与技术债

- `internal/response.Fail` 和 `Error` 使用 HTTP 200 搭配业务码 400/500。修改现有接口时先确认客户端契约，不做无说明批量切换。
- `service.UserServiceImpl.Login` 的历史明文密码升级忽略 Repository Update 错误。这是兼容路径，不是推荐吞错模式。
- 部分 Demo Handler 忽略参数转换或 Context 取值错误，不得复制到正式功能。
- `internal/redis.Init` 当前用 `%v` 包装 Ping 错误，无法通过 `errors.Is` 追踪；修改该路径时应改为 `%w` 并补测试。

## Review 检查

- 是否有 `_ = err`、空 catch 等价路径或只记录不返回的关键失败。
- 是否丢失 `%w` 导致上层无法分类。
- 是否在 Handler 返回了原始内部错误。
- 是否写响应后仍继续执行。
- 是否把网络超时、限流、参数错误和权限错误混成同一重试策略。
