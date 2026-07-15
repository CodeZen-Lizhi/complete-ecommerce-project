# 阶段 7 练习 3：限流与并发控制

## 练习目标

为接口、模型和 Embedding 增加限流、并发上限和取消支持。

## 前置知识

- Secret、PII、日志脱敏和最小权限。
- Logs、Metrics、Traces 与 SLI/SLO。
- 限流、健康检查、灰度、降级和回滚。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：定义请求级、用户级和下游级限额。
2. TODO 2：使用有界信号量控制并发，等待时响应 Context。
3. TODO 3：实现令牌桶或等价限速，并区分排队与拒绝。
4. TODO 4：确保释放许可，避免 panic/错误路径泄漏。
5. TODO 5：压测成功、超时、取消和突发流量，记录 P95 与拒绝率。

## 开始练习

```bash
go run ./examples/ai/phase7/03_rate_limit_and_concurrency
```

骨架固定为 `dry-run`，不会执行真实发布、故障注入或外部写操作。

## 验证方式

```bash
gofmt -w examples/ai/phase7/03_rate_limit_and_concurrency/*.go
go test -timeout=60s ./examples/ai/phase7/03_rate_limit_and_concurrency
go vet ./examples/ai/phase7/03_rate_limit_and_concurrency
```

## 完成标准

- 日志、指标和 Trace 不包含 Secret、PII 或高基数敏感标签。
- 超时、限流、健康和降级状态可以被测试和观察。
- 降级不伪造成功，恢复与回滚路径有明确验证。
- 部署模板不包含真实凭证，危险操作必须显式启用。

## 暂不实现

- 全局分布式配额平台。
