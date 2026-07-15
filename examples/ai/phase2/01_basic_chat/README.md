# 阶段 2 练习 1：基础模型 HTTP 调用

## 练习目标

直接使用 net/http 调用 OpenAI-compatible Chat Completions，掌握配置、请求、响应、状态码、Token 和超时。

## TODO 顺序

1. TODO 1：校验 API Key 占位符。
2. TODO 2：组装请求并序列化 JSON。
3. TODO 3：创建携带 Context 的 HTTP 请求。
4. TODO 4：限制并读取响应体，区分非 2xx。
5. TODO 5：解析回答与 Token，并处理空 choices。

## 开始练习

```bash
go run ./examples/ai/phase2/01_basic_chat
```

默认 API Key 或外部服务配置保持占位符时，程序应明确提示配置未完成，不得连接外部服务或产生付费请求。

## 验证方式

```bash
gofmt -w examples/ai/phase2/01_basic_chat/*.go
go test -timeout=60s ./examples/ai/phase2/01_basic_chat
go vet ./examples/ai/phase2/01_basic_chat
```

## 完成标准

- TODO 编号连续，错误使用 `%w` 保留底层原因。
- Context 超时和取消能够传播到模型或网络边界。
- 不打印 API Key、完整敏感 Prompt 或 Authorization。
- 删除核心实现后，测试或默认运行能够清楚暴露练习未完成。

## 暂不实现

- Eino 封装、流式输出、多轮历史和重试治理。
