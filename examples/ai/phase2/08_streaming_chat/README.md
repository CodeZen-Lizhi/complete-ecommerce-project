# 阶段 2 练习 8：ChatModel Stream

## 练习目标

使用 Eino BaseChatModel.Stream 消费增量输出，正确处理 Close、EOF、取消和中途错误。

## TODO 顺序

1. TODO 1：创建 ChatModel。
2. TODO 2：调用 Stream 并立即注册 Close。
3. TODO 3：循环 Recv，单独处理 io.EOF。
4. TODO 4：传播取消和中途错误。
5. TODO 5：输出非空 chunk 并拒绝空流。

## 开始练习

```bash
go run ./examples/ai/phase2/08_streaming_chat
```

默认 API Key 或外部服务配置保持占位符时，程序应明确提示配置未完成，不得连接外部服务或产生付费请求。

## 验证方式

```bash
gofmt -w examples/ai/phase2/08_streaming_chat/*.go
go test -timeout=60s ./examples/ai/phase2/08_streaming_chat
go vet ./examples/ai/phase2/08_streaming_chat
```

## 完成标准

- TODO 编号连续，错误使用 `%w` 保留底层原因。
- Context 超时和取消能够传播到模型或网络边界。
- 不打印 API Key、完整敏感 Prompt 或 Authorization。
- 删除核心实现后，测试或默认运行能够清楚暴露练习未完成。

## 暂不实现

- SSE HTTP 转发、断线续传和流式重试。
