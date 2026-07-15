# 阶段 1 练习 3：Redis Cache Aside、空值与失效

## 练习目标

为一个只读查询增加 Cache Aside，正确处理缓存命中、穿透防护、TTL 和写后失效。

## 实际修改位置建议

- `internal/redis/redis_client.go`
- `service/user_service.go`
- `repository/user_repository.go`

阶段 1 练习直接作用于现有电商项目；本目录只提供顺序、任务说明和验收清单，不复制一套业务代码。

## TODO 顺序

1. TODO 1：选择一个稳定的只读查询并定义带命名空间的缓存 Key。
2. TODO 2：先查 Redis，区分命中、未命中、空值标记和 Redis 故障。
3. TODO 3：未命中时查询 Repository，并把结果或短 TTL 空值标记写入缓存。
4. TODO 4：更新或删除成功后删除缓存，明确数据库成功但缓存删除失败的处理方式。
5. TODO 5：为正常命中、缓存穿透、Redis 故障和缓存失效补测试。

## 验证方式

```bash
go test -timeout=60s ./service ./internal/redis
go vet ./service ./internal/redis
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
