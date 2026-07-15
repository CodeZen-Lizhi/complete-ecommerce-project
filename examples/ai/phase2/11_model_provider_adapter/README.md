# 阶段 2 练习 11：模型厂商适配层

## 练习目标

让业务代码依赖项目自定义接口，通过配置切换具体模型实现。

## TODO 顺序

1. TODO 1：定义最小 chatProvider。
2. TODO 2：实现 Eino OpenAI-compatible adapter。
3. TODO 3：转换项目消息与 SDK 消息。
4. TODO 4：通过白名单 Factory 选择实现。
5. TODO 5：用 Fake 验证 Factory 和配置切换不影响调用方。
6. TODO 6：只通过项目接口发起调用，统一校验回答和 Token。

## 开始练习

```bash
go run ./examples/ai/phase2/11_model_provider_adapter
```

默认 API Key 或外部服务配置保持占位符时，程序应明确提示配置未完成，不得连接外部服务或产生付费请求。

## 验证方式

```bash
gofmt -w examples/ai/phase2/11_model_provider_adapter/*.go
go test -timeout=60s ./examples/ai/phase2/11_model_provider_adapter
go vet ./examples/ai/phase2/11_model_provider_adapter
```

## 完成标准

- TODO 编号连续，错误使用 `%w` 保留底层原因。
- Context 超时和取消能够传播到模型或网络边界。
- 不打印 API Key、完整敏感 Prompt 或 Authorization。
- 删除核心实现后，测试或默认运行能够清楚暴露练习未完成。

## 暂不实现

- 动态插件、跨厂商自动 fallback 和多区域路由。
