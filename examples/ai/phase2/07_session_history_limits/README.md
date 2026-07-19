# 阶段 2 练习 7：会话历史截断、过期与用户隔离

上一练习已经把完整消息历史保存到 Redis。本练习继续解决三个生产环境一定会遇到的问题：历史不能无限增长、长期不用的会话需要过期、同名 session 不能让不同用户互相读取。

## 练习目标

- 每次模型调用只携带最近 `maxHistoryTurns` 轮对话。
- 每轮成功后刷新 Redis TTL，让不活跃会话自动删除。
- Redis Key 同时包含 user ID 和 session ID，实现最小用户隔离。
- 模型失败时不保存半轮消息。

本练习选择“限制完整轮数”，暂不做精确 Token 计数。不同模型的 tokenizer 和上下文规则可能不同，先把消息边界、存储边界和隔离规则学清楚更重要。

## 关键设计

Redis List 只保存以下重复结构：

```text
User 1 -> Assistant 1 -> User 2 -> Assistant 2 -> ...
```

System Message 不写 Redis，而是在每次模型调用前重新注入：

```text
System -> 最近 N 轮 User/Assistant -> 当前 User
```

这样做有两个直接好处：

1. `LTRIM` 可以直接保留 List 尾部 `maxHistoryTurns * 2` 条消息，不会误删唯一的 System Message。
2. 每轮总是追加一对 User/Assistant，并按偶数条裁剪，不会主动制造半轮历史。

如果读取到奇数条消息，应把它视为数据损坏并明确报错，而不是把不完整上下文交给模型。

## TODO 顺序

1. TODO 1：校验 user ID、session ID 并生成隔离 Redis Key。
2. TODO 2：只读取最近 `maxHistoryTurns * 2` 条消息。
3. TODO 3：反序列化并拒绝奇数条损坏历史。
4. TODO 4：校验并序列化完整 User/Assistant 一轮。
5. TODO 5：在一次事务管道中追加、裁剪并刷新 TTL。
6. TODO 6：按当前用户和 session 重新加载历史。
7. TODO 7：按 System → 最近历史 → 当前 User 组装模型输入。
8. TODO 8：调用模型并提取 Assistant 文本，失败时不保存半轮。
9. TODO 9：模型成功后原子保存完整一轮。
10. TODO 10：演示同一用户的两轮关联对话。
11. TODO 11：使用另一用户和同名 session 验证隔离。

## 开始练习

1. 打开 `main.go`，从 `TODO 1` 开始按编号实现，不要跳步。
2. 不设置 `OPENAI_API_KEY`，先验证缺少配置时会安全退出：

   ```bash
   env -u OPENAI_API_KEY go run ./examples/ai/phase2/07_session_history_limits
   ```

   此时应只提示 `API Key 未配置`，不会连接 Redis，也不会调用模型。

3. 完成纯校验和 Redis 存储 TODO 后，确认本机 Redis 可用。
4. 仅在本地通过环境变量配置自己的 API Key，不要提交真实密钥：

   ```bash
   export OPENAI_API_KEY="你的本地测试 Key"
   ```

5. 运行：

   ```bash
   go run ./examples/ai/phase2/07_session_history_limits
   ```

## 如何验证三个目标

### 1. 历史截断

将演示扩展到超过 `maxHistoryTurns` 轮，再执行：

```bash
redis-cli LRANGE ai:phase2:user-session:user-1001:checkout-help 0 -1
```

List 元素数量不应超过 `maxHistoryTurns * 2`，且应始终是偶数。

### 2. 会话过期

执行：

```bash
redis-cli TTL ai:phase2:user-session:user-1001:checkout-help
```

成功保存一轮后 TTL 应接近配置值；再次成功对话后 TTL 会刷新。为了快速观察，可以把 `sessionTTL` 临时改成几十秒。

### 3. 用户隔离

程序会让 `user-1001` 和 `user-2002` 使用相同的 `checkout-help` session。Redis 中应出现两个不同 Key：

```text
ai:phase2:user-session:user-1001:checkout-help
ai:phase2:user-session:user-2002:checkout-help
```

`user-2002` 不应理解 `user-1001` 对“刚才类比”的追问，因为它没有访问前一个用户的历史。

## 完成后检查

```bash
gofmt -w examples/ai/phase2/07_session_history_limits/main.go
go test -timeout=60s ./examples/ai/phase2/07_session_history_limits
go vet ./examples/ai/phase2/07_session_history_limits
```

## 完成标准

- Redis Key 同时包含经过校验的 user ID 和 session ID。
- List 始终只保存完整 User/Assistant 对，数量为偶数且不超过历史上限。
- 每轮成功后刷新 TTL；模型失败时 Redis 不产生半轮写入。
- 两个用户使用同名 session 时不能读取彼此历史。

## 本练习暂时不做

- tokenizer 和精确 Token Budget。
- 旧历史摘要、长期记忆或向量检索。
- 并发写冲突、Lua 脚本和分布式锁。
- ChatTemplate、Stream、结构化输出、重试与限流。

它们会在后续练习中逐步加入。
