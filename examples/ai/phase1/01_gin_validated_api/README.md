# 阶段 1 练习 1：Gin 参数校验与统一响应

## 练习目标

在现有 Gin 路由中实现一个带明确输入边界、统一成功/失败响应和可测试状态码的接口。

## 实际修改位置建议

- `router/user_router.go`
- `handler/user_handler.go`
- `internal/response/result.go`

阶段 1 练习直接作用于现有电商项目；本目录只提供顺序、任务说明和验收清单，不复制一套业务代码。

## TODO 顺序

1. TODO 1：选择一个现有未完成接口，明确 Path、Method、输入字段和成功/失败状态码。
2. TODO 2：定义请求 DTO，并在 Handler 边界完成必填、格式、范围和空白字符校验。
3. TODO 3：调用 Service 接口，不在 Handler 中写业务规则或访问数据库。
4. TODO 4：使用项目统一响应结构映射成功、业务错误和内部错误。
5. TODO 5：使用 httptest 覆盖合法输入、缺失字段、非法字段和 Service 失败。

## 验证方式

```bash
go test -timeout=60s ./handler ./router/...
go vet ./handler ./router/...
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
