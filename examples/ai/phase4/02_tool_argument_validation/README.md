# 阶段 4 练习 2：Tool 参数严格校验

## 练习目标

把模型生成的参数视为不可信输入，校验必填、类型、范围和未知字段。

## 前置知识

- 模型只生成 Tool Call，应用负责校验、授权和执行。
- 模型参数是不可信输入，Context 中的身份才是可信边界。
- 写操作需要确认、幂等、审计和明确失败处理。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：使用严格 JSON Decoder 拒绝未知字段。
2. TODO 2：校验必填字符串、枚举、数值范围和数组上限。
3. TODO 3：区分语法错误、Schema 错误和业务校验错误。
4. TODO 4：生成对模型友好的有限错误信息，不泄露内部实现。
5. TODO 5：用表驱动测试覆盖边界和恶意输入。

## 开始练习

```bash
go run ./examples/ai/phase4/02_tool_argument_validation
```

骨架默认使用本地接口和安全配置，不连接生产数据库，也不执行真实写操作。

## 验证方式

```bash
gofmt -w examples/ai/phase4/02_tool_argument_validation/*.go
go test -timeout=60s ./examples/ai/phase4/02_tool_argument_validation
go vet ./examples/ai/phase4/02_tool_argument_validation
```

## 完成标准

- 参数、身份、租户、权限和结果上限都在工具边界验证。
- Context 取消和依赖错误清楚传播，错误不泄露敏感内部信息。
- 只读调用可测试，写操作默认禁用并具备确认、幂等和审计。
- 不让模型自由决定可信身份或绕过固定业务规则。

## 暂不实现

- 动态 Schema 和自动纠错循环。
