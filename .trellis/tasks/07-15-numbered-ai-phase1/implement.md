# 阶段 1：Go 后端基础编号练习执行计划

## 实现清单

- [x] 生成 `01_gin_validated_api`：Gin 参数校验与统一响应。
- [x] 生成 `02_gorm_crud_pagination_transaction`：GORM CRUD、分页、事务与 Explain。
- [x] 生成 `03_redis_cache_aside`：Redis Cache Aside、空值与失效。
- [x] 生成 `04_jwt_auth_and_authorization`：JWT 登录、鉴权与后端授权。
- [x] 生成 `05_http_client_resilience`：HTTP Client 超时与非 2xx 处理。
- [x] 生成 `06_worker_pool`：Worker Pool、取消与错误收集。
- [x] 生成 `07_backend_tests`：Service、Handler、数据库与失败路径测试。

## 一致性检查

- [x] 验证目录编号从 `01` 连续到 `07`。
- [x] 验证每个目录包含 README，章节和 TODO 编号连续。
- [x] 验证默认配置不会连接外部服务或产生付费请求。
- [x] 搜索本阶段旧路径、重复目录、真实密钥和未处理错误。

## 验证

```bash
gofmt -w <阶段 1 新增或迁移的 Go 文件>
go test -timeout=60s ./examples/ai/phase1/...
go vet ./examples/ai/phase1/...
git diff --check
```

阶段 1 完成后执行 `trellis-check` 和 `go-review`；涉及 SQL/Schema/查询时追加 `sql-code-review`。
