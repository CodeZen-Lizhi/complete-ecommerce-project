# 阶段 2 练习 5：Redis 会话历史

## 练习目标

使用 Redis List 按 session 保存消息，每一轮都从 Redis 重新加载并组装模型上下文。

## TODO 顺序

1. TODO 1：校验 session ID 并生成命名空间 Key。
2. TODO 2：读取并反序列化完整历史。
3. TODO 3：序列化并按顺序追加本轮消息。
4. TODO 4：按 System/历史/当前 User 组装输入。
5. TODO 5：模型成功后保存新增消息。

## 开始练习

```bash
go run ./examples/ai/phase2/05_redis_session_history
```

默认 API Key 或外部服务配置保持占位符时，程序应明确提示配置未完成，不得连接外部服务或产生付费请求。

## 验证方式

```bash
gofmt -w examples/ai/phase2/05_redis_session_history/*.go
go test -timeout=60s ./examples/ai/phase2/05_redis_session_history
go vet ./examples/ai/phase2/05_redis_session_history
```

## 完成标准

- TODO 编号连续，错误使用 `%w` 保留底层原因。
- Context 超时和取消能够传播到模型或网络边界。
- 不打印 API Key、完整敏感 Prompt 或 Authorization。
- 删除核心实现后，测试或默认运行能够清楚暴露练习未完成。

## 暂不实现

- 用户隔离、TTL、历史截断和并发冲突治理。
