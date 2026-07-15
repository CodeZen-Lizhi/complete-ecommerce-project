# 阶段 4 练习 3：Context 身份与再次授权

## 练习目标

把用户和租户身份沿 Context 传递，并在工具执行边界再次授权。

## 前置知识

- 模型只生成 Tool Call，应用负责校验、授权和执行。
- 模型参数是不可信输入，Context 中的身份才是可信边界。
- 写操作需要确认、幂等、审计和明确失败处理。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：定义可信 Identity 并由应用边界注入 Context。
2. TODO 2：工具只从 Context 读取身份，不接受模型参数中的 user ID 作为可信值。
3. TODO 3：校验租户、资源归属和权限动作。
4. TODO 4：区分未认证、无权限和资源不存在。
5. TODO 5：记录审计字段但不记录 Token、Cookie 或完整参数。

## 开始练习

```bash
go run ./examples/ai/phase4/03_context_identity_authorization
```

骨架默认使用本地接口和安全配置，不连接生产数据库，也不执行真实写操作。

## 验证方式

```bash
gofmt -w examples/ai/phase4/03_context_identity_authorization/*.go
go test -timeout=60s ./examples/ai/phase4/03_context_identity_authorization
go vet ./examples/ai/phase4/03_context_identity_authorization
```

## 完成标准

- 参数、身份、租户、权限和结果上限都在工具边界验证。
- Context 取消和依赖错误清楚传播，错误不泄露敏感内部信息。
- 只读调用可测试，写操作默认禁用并具备确认、幂等和审计。
- 不让模型自由决定可信身份或绕过固定业务规则。

## 暂不实现

- 完整 IAM 平台和跨租户委托。
