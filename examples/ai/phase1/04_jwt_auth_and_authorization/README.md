# 阶段 1 练习 4：JWT 登录、鉴权与后端授权

## 练习目标

区分身份认证与业务授权，完成安全的 JWT 生成、解析、中间件注入和可信权限判断。

## 实际修改位置建议

- `util/jwtutil.go`
- `middleware/middleware.go`
- `service/user_service.go`

阶段 1 练习直接作用于现有电商项目；本目录只提供顺序、任务说明和验收清单，不复制一套业务代码。

## TODO 顺序

1. TODO 1：检查 JWT issuer、audience、过期时间和签名算法约束。
2. TODO 2：登录成功后生成 Token，不在日志或响应外字段泄露敏感信息。
3. TODO 3：中间件从 Authorization Bearer 读取并验证 Token，再把可信用户身份写入 Context。
4. TODO 4：在 Service 或授权边界根据可信身份判断资源归属，不信任客户端 Role Header。
5. TODO 5：覆盖合法、过期、签名算法错误、issuer/audience 错误和越权访问测试。

## 验证方式

```bash
go test -timeout=60s ./util ./middleware ./service
go vet ./util ./middleware ./service
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
