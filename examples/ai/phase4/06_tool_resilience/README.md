# 阶段 4 练习 6：工具失败治理

## 练习目标

处理超时、下游失败、返回过大和有限重试。

## 前置知识

- 模型只生成 Tool Call，应用负责校验、授权和执行。
- 模型参数是不可信输入，Context 中的身份才是可信边界。
- 写操作需要确认、幂等、审计和明确失败处理。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：为每次工具调用设置独立超时。
2. TODO 2：分类取消、超时、临时下游错误和永久业务错误。
3. TODO 3：只对幂等只读调用做有限重试。
4. TODO 4：限制响应字节数和列表条数。
5. TODO 5：记录耗时、重试次数和失败类型，并验证 goroutine 不泄漏。

## 开始练习

```bash
go run ./examples/ai/phase4/06_tool_resilience
```

骨架默认使用本地接口和安全配置，不连接生产数据库，也不执行真实写操作。

## 验证方式

```bash
gofmt -w examples/ai/phase4/06_tool_resilience/*.go
go test -timeout=60s ./examples/ai/phase4/06_tool_resilience
go vet ./examples/ai/phase4/06_tool_resilience
```

## 完成标准

- 参数、身份、租户、权限和结果上限都在工具边界验证。
- Context 取消和依赖错误清楚传播，错误不泄露敏感内部信息。
- 只读调用可测试，写操作默认禁用并具备确认、幂等和审计。
- 不让模型自由决定可信身份或绕过固定业务规则。

## 暂不实现

- 跨服务熔断平台和分布式限流。
