# 数据库规范

## 技术与入口

项目使用 MySQL、GORM 和 `go-sql-driver/mysql`。连接初始化在 `internal/mysql/mysql.go`，全局 `mysql.DB` 只用于启动组装和事务入口；业务查询通过注入的 Repository 执行。

连接池参数来自 `configs/config.<env>.yaml`：

- `max_open_conns`
- `max_idle_conns`
- `conn_max_lifetime`（代码中按秒转换）
- `debug`（决定 GORM 日志级别）

不要在 Handler 或 Service 中再次调用 `gorm.Open`。

## 模型约定

持久化结构位于 `model/`，使用 GORM 标签显式描述主键、列名、大小、默认值和索引。每个模型通过 `TableName()` 绑定真实表名：

- `model.User` -> `users`
- `model.Dict` -> `s_dict`
- `model.DictItem` -> `s_dict_item`

ID 策略并不统一：User 在 `BeforeCreate` 中调用 `util.GenID()`，Dict/DictItem 标签使用自增。新增模型必须依据真实数据库设计选择策略，不能假设全项目都使用雪花 ID 或自增 ID。

可空数据库字段当前使用指针类型，例如 `*string`、`*int64`、`*time.Time`。时间和软删除字段采用现有 `create_time`、`update_time`、`del_flag` 命名。

## 查询规范

- 用户输入只能作为 GORM 参数值传入，例如 `Where("email = ? and del_flag = ?", email, "0")`。
- 动态表名、列名和排序字段不能直接拼接用户输入；必须使用白名单。
- 当前软删除是业务字段 `del_flag`，不是 `gorm.DeletedAt`。需要排除删除数据时显式加 `del_flag = "0"`。
- `gorm.ErrRecordNotFound` 用 `errors.Is` 判断；“不存在是否为错误”由 Repository 方法契约决定。
- 列表查询必须有分页、稳定排序和合理上限，不能一次加载无界数据。
- 避免循环查库和逐条更新；优先批量查询/写入，并验证生成 SQL。

`UserRepository.FindByID` 当前没有过滤 `del_flag`，这是现有不一致，不应复制到新的用户查询方法。

## 更新语义

`UserRepository.Update` 当前调用 `db.Updates(user)`。GORM 对 struct 的 `Updates` 默认忽略零值，因此当业务需要把字段更新为 `0`、`false` 或空字符串时，应显式使用 `Select` 或受控的 `map[string]any`，并写测试证明字段集合。

不要直接对未经白名单筛选的请求 DTO 执行全字段更新。密码、审计字段、主键和删除标记必须由业务规则控制。

## 事务

统一使用 `internal/mysql.Transaction(ctx, fn)`：

- `fn == nil` 和 DB 未初始化会返回明确错误。
- 调用链 ctx 通过 `DB.WithContext(ctx)` 进入事务。
- Service 在事务回调中调用 `repo.WithDB(tx)`，后续全部操作使用该 Repository。
- 回调返回非 nil error 触发回滚，返回 nil 提交。

参考：`service.UserServiceImpl.Register`。不要混用事务外 Repository，否则部分写入不会回滚。

## Schema 与迁移

仓库当前没有 SQL migration 目录，也没有 `AutoMigrate` 调用。`model.DictItem.Indexes()` 只是普通 Go 方法，不能视为数据库迁移机制或已落库约束。

因此：

- 不凭 GORM 标签推断线上 Schema 已同步。
- Schema、索引或字段修改前必须先确认数据库事实和发布方式。
- DB 工具不可用且仓库没有迁移事实时，默认不创建 Schema 变更。
- 新增索引必须同时说明目标查询、字段顺序、唯一性、回滚和验证 SQL。

## 验证

- Repository 单元测试至少覆盖成功、未找到、数据库错误和软删除条件。
- 事务测试覆盖提交与回滚，确认所有操作使用同一个 tx。
- 性能敏感查询记录 `EXPLAIN`、行数和索引命中；列表接口同时检查分页上限。
- 涉及 SQL、Schema、索引、权限或一致性的修改，除 Go review 外追加 SQL 专项 review。
