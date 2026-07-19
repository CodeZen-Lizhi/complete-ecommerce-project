# 阶段 7 练习 7：健康检查、灰度与回滚

## 练习目标

配置健康检查、灰度发布和可验证回滚。

## 前置知识

- Secret、PII、日志脱敏和最小权限。
- Logs、Metrics、Traces 与 SLI/SLO。
- 限流、健康检查、灰度、降级和回滚。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：区分 liveness、readiness 和依赖降级状态。
2. TODO 2：定义版本、配置和模型路由的发布标识。
3. TODO 3：让少量测试流量进入新版本并比较错误率、P95 和质量。
4. TODO 4：达到阈值时停止扩量并回滚。
5. TODO 5：验证回滚不破坏会话、任务和索引兼容性。

## 开始练习

先修改 `exercise.go` 顶部的 `releaseControlBaseURL` 和 `releaseControlAPIKey`，并确保地址只指向练习环境。

```bash
go run ./examples/ai/phase7/07_health_gray_release_rollback
```

完成 TODO 后需要连接练习专用发布控制端点，执行小流量扩量和回滚演练；禁止指向生产环境。

## 验证方式

```bash
gofmt -w examples/ai/phase7/07_health_gray_release_rollback/*.go
go test -timeout=60s ./examples/ai/phase7/07_health_gray_release_rollback
go vet ./examples/ai/phase7/07_health_gray_release_rollback
```

## 完成标准

- 使用顶部 `releaseControlBaseURL` 连接真实练习环境并执行权重调整与回滚；只计算枚举结果不算完成。

- 日志、指标和 Trace 不包含 Secret、PII 或高基数敏感标签。
- 超时、限流、健康和降级状态可以被测试和观察。
- 降级不伪造成功，恢复与回滚路径有明确验证。
- 部署模板不包含真实凭证，危险操作必须显式启用。

## 暂不实现

- 真实生产流量切换。
