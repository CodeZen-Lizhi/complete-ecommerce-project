# 阶段 2 练习 2：Eino ChatModel Generate

## 练习目标

使用 Eino OpenAI-compatible ChatModel 完成与上一题相同的阻塞式调用。

## TODO 顺序

1. TODO 1：校验模型配置。
2. TODO 2：创建 openai.ChatModelConfig 和 ChatModel。
3. TODO 3：使用 schema.Message 组装 System/User 消息。
4. TODO 4：调用 Generate 并处理 nil 响应。

## 开始练习

```bash
go run ./examples/ai/phase2/02_eino_generate
```

默认 API Key 或外部服务配置保持占位符时，程序应明确提示配置未完成，不得连接外部服务或产生付费请求。

## 验证方式

```bash
gofmt -w examples/ai/phase2/02_eino_generate/*.go
go test -timeout=60s ./examples/ai/phase2/02_eino_generate
go vet ./examples/ai/phase2/02_eino_generate
```

## 完成标准

- TODO 编号连续，错误使用 `%w` 保留底层原因。
- Context 超时和取消能够传播到模型或网络边界。
- 不打印 API Key、完整敏感 Prompt 或 Authorization。
- 删除核心实现后，测试或默认运行能够清楚暴露练习未完成。

## 暂不实现

- 多轮历史、Stream、Tool Calling 和模型适配层。
