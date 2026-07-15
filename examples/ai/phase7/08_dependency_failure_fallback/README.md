# 阶段 7 练习 8：依赖故障与降级

## 练习目标

模拟主模型或向量库不可用，并验证受控降级。

## 前置知识

- Secret、PII、日志脱敏和最小权限。
- Logs、Metrics、Traces 与 SLI/SLO。
- 限流、健康检查、灰度、降级和回滚。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：定义可降级与不可降级场景，禁止假成功。
2. TODO 2：注入模型超时、429、5xx 和向量库连接失败。
3. TODO 3：只对安全只读请求启用备用模型或缓存结果。
4. TODO 4：降级响应明确标记能力限制并记录原因。
5. TODO 5：验证恢复后自动退出降级，且不会形成重试风暴。

## 开始练习

```bash
go run ./examples/ai/phase7/08_dependency_failure_fallback
```

骨架固定为 `dry-run`，不会执行真实发布、故障注入或外部写操作。

## 验证方式

```bash
gofmt -w examples/ai/phase7/08_dependency_failure_fallback/*.go
go test -timeout=60s ./examples/ai/phase7/08_dependency_failure_fallback
go vet ./examples/ai/phase7/08_dependency_failure_fallback
```

## 完成标准

- 日志、指标和 Trace 不包含 Secret、PII 或高基数敏感标签。
- 超时、限流、健康和降级状态可以被测试和观察。
- 降级不伪造成功，恢复与回滚路径有明确验证。
- 部署模板不包含真实凭证，危险操作必须显式启用。

## 暂不实现

- 跨区域容灾和多云路由。
