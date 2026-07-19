# 阶段 2 练习 5：程序内多轮对话

## 练习目标

在进程内保存完整 Assistant 消息，并在第二轮重新发送 System、User 和 Assistant 历史。

## TODO 顺序

1. TODO 1：初始化 System 消息。
2. TODO 2：追加第一轮 User 并调用模型。
3. TODO 3：保存完整第一轮 Assistant 消息。
4. TODO 4：追加第二轮 User 并携带全部历史调用。
5. TODO 5：提取并展示 AssistantGenText。

## 开始练习

```bash
go run ./examples/ai/phase2/05_in_memory_multi_turn
```

默认 API Key 或外部服务配置保持占位符时，程序应明确提示配置未完成，不得连接外部服务或产生付费请求。

## 验证方式

```bash
gofmt -w examples/ai/phase2/05_in_memory_multi_turn/*.go
go test -timeout=60s ./examples/ai/phase2/05_in_memory_multi_turn
go vet ./examples/ai/phase2/05_in_memory_multi_turn
```

## 完成标准

- TODO 编号连续，错误使用 `%w` 保留底层原因。
- Context 超时和取消能够传播到模型或网络边界。
- 不打印 API Key、完整敏感 Prompt 或 Authorization。
- 删除核心实现后，测试或默认运行能够清楚暴露练习未完成。

## 暂不实现

- Redis 持久化、历史截断、TTL 和并发写治理。
