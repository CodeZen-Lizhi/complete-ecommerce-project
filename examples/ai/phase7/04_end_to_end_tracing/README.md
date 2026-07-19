# 阶段 7 练习 4：端到端 Trace

## 练习目标

使用 trace ID 串联模型、Retriever、Rerank、Tool 和 Agent。

## 前置知识

- Secret、PII、日志脱敏和最小权限。
- Logs、Metrics、Traces 与 SLI/SLO。
- 限流、健康检查、灰度、降级和回滚。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：创建真实 OpenTelemetry stdout Exporter、SDK TracerProvider 和入口根 Span。
2. TODO 2：使用 `otel.Tracer` 为各组件创建有父子关系的 Span。
3. TODO 3：记录耗时、状态、模型名和结果数量等低敏属性。
4. TODO 4：错误写入 span，但不记录 Prompt、Secret 或完整文档。
5. TODO 5：验证一次请求能够还原完整调用链和失败节点。

## 开始练习

```bash
go run ./examples/ai/phase7/04_end_to_end_tracing
```

完成 TODO 后程序会把真实 OpenTelemetry Span 导出到标准输出，可直接检查父子关系和错误状态。

## 验证方式

```bash
gofmt -w examples/ai/phase7/04_end_to_end_tracing/*.go
go test -timeout=60s ./examples/ai/phase7/04_end_to_end_tracing
go vet ./examples/ai/phase7/04_end_to_end_tracing
```

## 完成标准

- 实际创建 `sdktrace.TracerProvider` 并导出 Span；只实现本地 `tracer` 接口不算完成。

- 日志、指标和 Trace 不包含 Secret、PII 或高基数敏感标签。
- 超时、限流、健康和降级状态可以被测试和观察。
- 降级不伪造成功，恢复与回滚路径有明确验证。
- 部署模板不包含真实凭证，危险操作必须显式启用。

## 暂不实现

- 供应商专属 APM 平台配置。
