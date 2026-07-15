# 阶段 1 练习 5：HTTP Client 超时与非 2xx 处理

## 练习目标

封装一个可测试的外部 HTTP 调用，正确传播 Context、限制响应体并区分网络和 HTTP 错误。

## 实际修改位置建议

- `建议新增独立 client 包或放入明确的基础设施目录`
- `调用方 Service`

阶段 1 练习直接作用于现有电商项目；本目录只提供顺序、任务说明和验收清单，不复制一套业务代码。

## TODO 顺序

1. TODO 1：定义最小 Client 接口、请求/响应 DTO 和可配置 Base URL。
2. TODO 2：复用带超时的 http.Client，并使用 NewRequestWithContext 传播取消。
3. TODO 3：设置必要 Header，发送请求后立即关闭响应体。
4. TODO 4：限制响应体大小，区分非 2xx、超时、取消、解码失败和空响应。
5. TODO 5：使用 httptest.Server 覆盖成功、慢响应、非 2xx、超大响应和非法 JSON。

## 验证方式

```bash
go test -timeout=60s ./...
go vet ./...
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
