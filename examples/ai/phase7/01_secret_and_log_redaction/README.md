# 阶段 7 练习 1：Secret 管理与日志脱敏

## 练习目标

使用环境变量或 Secret 管理 Key，并确保日志隐藏敏感信息。

## 前置知识

- Secret、PII、日志脱敏和最小权限。
- Logs、Metrics、Traces 与 SLI/SLO。
- 限流、健康检查、灰度、降级和回滚。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：从环境读取配置，拒绝源码中的真实默认密钥。
2. TODO 2：区分必需配置、可选配置和安全占位符。
3. TODO 3：实现结构化日志字段白名单。
4. TODO 4：对 Authorization、Cookie、Token、密码和 PII 做脱敏。
5. TODO 5：用测试验证错误和日志中不包含原始 Secret。

## 开始练习

```bash
go run ./examples/ai/phase7/01_secret_and_log_redaction
```

骨架固定为 `dry-run`，不会执行真实发布、故障注入或外部写操作。

## 验证方式

```bash
gofmt -w examples/ai/phase7/01_secret_and_log_redaction/*.go
go test -timeout=60s ./examples/ai/phase7/01_secret_and_log_redaction
go vet ./examples/ai/phase7/01_secret_and_log_redaction
```

## 完成标准

- 日志、指标和 Trace 不包含 Secret、PII 或高基数敏感标签。
- 超时、限流、健康和降级状态可以被测试和观察。
- 降级不伪造成功，恢复与回滚路径有明确验证。
- 部署模板不包含真实凭证，危险操作必须显式启用。

## 暂不实现

- 云厂商 Secret Manager 和自动轮换。
