# 阶段 4 练习 7：安全写操作模拟

## 练习目标

模拟取消订单/退款工具，加入确认、幂等和审计。

## 前置知识

- 模型只生成 Tool Call，应用负责校验、授权和执行。
- 模型参数是不可信输入，Context 中的身份才是可信边界。
- 写操作需要确认、幂等、审计和明确失败处理。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：把写操作默认设置为 dry-run，先返回影响摘要。
2. TODO 2：要求应用侧确认令牌，模型文本不能代替用户确认。
3. TODO 3：使用幂等键避免重复取消或退款。
4. TODO 4：再次检查身份、订单状态和金额边界。
5. TODO 5：记录审计事件并处理部分失败，测试重复调用不产生第二次副作用。

## 开始练习

```bash
go run ./examples/ai/phase4/07_safe_write_tool
```

骨架默认使用本地接口和安全配置，不连接生产数据库，也不执行真实写操作。

## 验证方式

```bash
gofmt -w examples/ai/phase4/07_safe_write_tool/*.go
go test -timeout=60s ./examples/ai/phase4/07_safe_write_tool
go vet ./examples/ai/phase4/07_safe_write_tool
```

## 完成标准

- 参数、身份、租户、权限和结果上限都在工具边界验证。
- Context 取消和依赖错误清楚传播，错误不泄露敏感内部信息。
- 只读调用可测试，写操作默认禁用并具备确认、幂等和审计。
- 不让模型自由决定可信身份或绕过固定业务规则。

## 暂不实现

- 真实支付、退款和生产订单变更。
