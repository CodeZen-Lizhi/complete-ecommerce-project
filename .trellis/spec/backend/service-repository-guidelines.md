# Service 与 Repository 规范

## Service 职责

Service 负责业务规则、跨 Repository 协作、密码处理和事务边界。参考 `service/user_service.go`：

- 接口定义调用契约，例如 `UserService`。
- 实现通过构造函数注入 Repository，不直接读取全局 Gin Context。
- 使用 `var _ UserService = (*UserServiceImpl)(nil)` 做编译期接口检查。
- 公开方法将 `context.Context` 作为第一个参数。
- 可供上层稳定判断的业务状态使用包级哨兵错误，例如 `ErrUserAlreadyExists`。

Service 不写 HTTP 响应，也不决定具体 HTTP 状态码。

## Repository 职责

Repository 只负责数据访问和持久化错误传播。参考 `repository/user_repository.go`：

- 接口与 GORM 实现位于同一业务仓储文件。
- 构造函数显式接收 `*gorm.DB`，例如 `NewUserRepository(db)`。
- 实现保存注入的 `db`，不在每个方法重新打开连接。
- 事务内通过 `WithDB(tx)` 返回绑定到事务 DB 的新实例，不能修改共享 Repository 的 DB 指针。
- 查询使用 GORM 参数绑定，例如 `Where("username = ? and del_flag = ?", username, "0")`。
- 数据库错误通过 `%w` 保留错误链，供 Service 使用 `errors.Is` 判断 `gorm.ErrRecordNotFound`。

## 事务协作

`service.UserServiceImpl.Register` 是当前事务范例：

```text
Service 接收 ctx
-> mysql.Transaction(ctx, fn)
-> repo.WithDB(tx)
-> 同一 repo 完成检查和写入
-> fn 返回 error 决定提交或回滚
```

- 事务边界放在 Service，而不是 Handler。
- 事务中的全部数据库操作必须使用 `tx` 派生的 Repository。
- 不在事务回调中启动脱离事务生命周期的 goroutine。
- 业务错误直接返回，让 GORM 回滚；不要在回调中吞错后返回 `nil`。

## Context 传播

现有 `UserRepository` 方法尚未接收 `context.Context`，因此 `FindByID` 等 Service 方法虽接收 ctx，却没有传到 GORM，这是已知技术债。新增 Repository API 应采用：

```go
FindByID(ctx context.Context, id int64) (*model.User, error)
```

实现使用 `r.db.WithContext(ctx)`。修改既有接口时必须同步所有实现、调用方和测试，不要只在 Service 签名上保留一个未使用的 ctx。

## 业务数据与副作用

- 密码哈希属于 Service，当前使用 bcrypt，参考 `hashPassword`。
- 不向 Handler 返回带密码的持久化模型；为对外响应定义不含敏感字段的 DTO。
- 默认不要修改调用方传入的对象。当前注册流程原地替换 `user.Password` 是既有行为；新流程优先复制必要字段再持久化。
- `Login` 中忽略历史密码升级写库错误是兼容代码，不得作为吞错范例。新副作用失败要明确记录、返回或进入可观测补偿路径。

## 依赖组装

生产依赖只在 `container.GetInstance` 中组装。新增 Service 后：

1. 给 `Container` 增加接口类型字段。
2. 使用已初始化的 `*gorm.DB` 构造 Repository。
3. 用 Repository 构造 Service。
4. 通过 `InjectContainerMiddleware` 注入，不在 Handler 内自行 `new`。

测试可直接构造 Service 并注入 Mock/Fake Repository，不依赖全局 Container。

## 禁止事项

- Handler 直接调用 GORM、MySQL 或 Redis完成业务写入。
- Repository 返回面向用户的中文文案或 Gin 响应。
- Service 依赖 `*gin.Context`。
- 在循环里无界逐条查库或调用远端服务。
- 用字符串匹配判断可包装错误；应使用 `errors.Is` / `errors.As`。
