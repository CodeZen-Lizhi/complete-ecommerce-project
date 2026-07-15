# 阶段 7 练习 2：Prompt Injection 防御

## 练习目标

构造直接和间接 Prompt Injection，验证 RAG 与 Tool 权限边界。

## 前置知识

- Secret、PII、日志脱敏和最小权限。
- Logs、Metrics、Traces 与 SLI/SLO。
- 限流、健康检查、灰度、降级和回滚。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：准备正常、直接注入和恶意文档三类 fixture。
2. TODO 2：把检索内容标记为不可信数据，不让其覆盖 System 规则。
3. TODO 3：在 Retriever 和 Tool 执行阶段强制租户与权限过滤。
4. TODO 4：拒绝模型请求泄露 Prompt、Secret 或越权调用。
5. TODO 5：记录攻击类型和阻断原因，并保留人工复核证据。

## 开始练习

```bash
go run ./examples/ai/phase7/02_prompt_injection_defense
```

骨架固定为 `dry-run`，不会执行真实发布、故障注入或外部写操作。

## 验证方式

```bash
gofmt -w examples/ai/phase7/02_prompt_injection_defense/*.go
go test -timeout=60s ./examples/ai/phase7/02_prompt_injection_defense
go vet ./examples/ai/phase7/02_prompt_injection_defense
```

## 完成标准

- 日志、指标和 Trace 不包含 Secret、PII 或高基数敏感标签。
- 超时、限流、健康和降级状态可以被测试和观察。
- 降级不伪造成功，恢复与回滚路径有明确验证。
- 部署模板不包含真实凭证，危险操作必须显式启用。

## 暂不实现

- 宣称完全消除 Prompt Injection。
