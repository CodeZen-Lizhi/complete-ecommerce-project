# 阶段 7 练习 5：Metrics 与告警

## 练习目标

记录耗时、Token、成本、错误率和队列长度并配置告警规则。

## 前置知识

- Secret、PII、日志脱敏和最小权限。
- Logs、Metrics、Traces 与 SLI/SLO。
- 限流、健康检查、灰度、降级和回滚。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：定义低基数指标名称、标签和单位。
2. TODO 2：记录请求量、成功率、P50/P95、Token、成本和队列长度。
3. TODO 3：避免把 user ID、session ID 或 Prompt 放入标签。
4. TODO 4：为错误率、延迟、队列积压和预算设置阈值。
5. TODO 5：用本地 fixture 模拟告警触发与恢复。

## 开始练习

```bash
go run ./examples/ai/phase7/05_metrics_and_alerts
```

骨架固定为 `dry-run`，不会执行真实发布、故障注入或外部写操作。

## 验证方式

```bash
gofmt -w examples/ai/phase7/05_metrics_and_alerts/*.go
go test -timeout=60s ./examples/ai/phase7/05_metrics_and_alerts
go vet ./examples/ai/phase7/05_metrics_and_alerts
```

## 完成标准

- 日志、指标和 Trace 不包含 Secret、PII 或高基数敏感标签。
- 超时、限流、健康和降级状态可以被测试和观察。
- 降级不伪造成功，恢复与回滚路径有明确验证。
- 部署模板不包含真实凭证，危险操作必须显式启用。

## 暂不实现

- 真实通知渠道和 24x7 值班流程。
