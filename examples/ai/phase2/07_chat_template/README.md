# 阶段 2 练习 7：Eino ChatTemplate 与变量

## 练习目标

使用 prompt.FromMessages、MessagesPlaceholder 和 Format 管理 System、历史与当前问题。

## TODO 顺序

1. TODO 1：创建包含 System、可选历史占位符和 User 的 FString ChatTemplate。
2. TODO 2：调用模板构造函数并保留错误链。
3. TODO 3：构造与模板键完全一致的变量 map。
4. TODO 4：调用 Format 并处理缺失变量或格式错误。
5. TODO 5：校验消息顺序为 System → 可选历史 → User。

## 开始练习

```bash
go run ./examples/ai/phase2/07_chat_template
```

默认 API Key 或外部服务配置保持占位符时，程序应明确提示配置未完成，不得连接外部服务或产生付费请求。

## 验证方式

```bash
gofmt -w examples/ai/phase2/07_chat_template/*.go
go test -timeout=60s ./examples/ai/phase2/07_chat_template
go vet ./examples/ai/phase2/07_chat_template
```

## 完成标准

- TODO 编号连续，错误使用 `%w` 保留底层原因。
- Context 超时和取消能够传播到模型或网络边界。
- 不打印 API Key、完整敏感 Prompt 或 Authorization。
- 删除核心实现后，测试或默认运行能够清楚暴露练习未完成。

## 暂不实现

- 模型调用、Stream 和复杂 Jinja2 模板。
