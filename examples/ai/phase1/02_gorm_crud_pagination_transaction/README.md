# 阶段 1 练习 2：GORM CRUD、分页、事务与 Explain

## 练习目标

在现有 Repository/Service 分层中完成可控的 CRUD、分页、事务和慢查询分析。

## 实际修改位置建议

- `repository/user_repository.go`
- `service/user_service.go`
- `internal/mysql/tx.go`

阶段 1 练习直接作用于现有电商项目；本目录只提供顺序、任务说明和验收清单，不复制一套业务代码。

## TODO 顺序

1. TODO 1：确认真实模型字段和现有 Repository 契约，不根据 README 猜测 Schema。
2. TODO 2：实现带 context 的创建、单条查询、更新和软删除路径。
3. TODO 3：实现有最大页大小的分页查询，并返回总数与当前页数据。
4. TODO 4：在 Service 中使用同一个事务派生 Repository 完成多步写入和失败回滚。
5. TODO 5：对关键查询记录 DryRun SQL 或真实环境 EXPLAIN，并说明索引命中情况。
6. TODO 6：补充 NotFound、更新零值、分页边界和事务回滚测试。

## 验证方式

```bash
go test -timeout=60s ./repository ./service ./internal/mysql
go vet ./repository ./service ./internal/mysql
git diff --check
```

## 完成标准

- 改动遵循 Handler → Service → Repository 依赖方向。
- 错误和取消能够清楚传播，没有静默成功或吞错。
- 测试覆盖本练习的成功路径和主要失败路径。
- 没有修改无关模块，也没有提交真实密钥或本地配置。

## 暂不实现

- 与本练习无关的全项目重构。
- 生产环境迁移、真实数据清理或大规模性能压测。
