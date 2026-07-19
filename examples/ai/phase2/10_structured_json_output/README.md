# 阶段 2 练习 10：结构化 JSON 输出

## 练习目标

完成 Prompt 约束、原生 JSON Schema Structured Outputs 和两种模式的失败边界对比。

## 9A：Prompt 约束与应用端校验

1. TODO 1：编写禁止 Markdown 代码块和额外文本的 System Prompt。
2. TODO 2：在 User Message 中提供固定 JSON 结构和当前问题。
3. TODO 3：通过 Model Client 生成原始文本，模型错误不得当成空 JSON。
4. TODO 4：把原始文本直接交给严格解析器，不静默截取或修复。
5. TODO 5：使用 Decoder 严格解析并拒绝未知字段。
6. TODO 6：拒绝第二个 JSON 值和尾随非空内容。
7. TODO 7：校验必填、数组元素和 Confidence 范围。
8. TODO 8：区分语法、结构和业务字段错误，并保留错误链。

## 9B：原生 JSON Schema Structured Outputs

9. TODO 9：从 `structuredAnswer` 生成 JSON Schema。
10. TODO 10：通过 Eino `ResponseFormat` 传入 `json_schema`，并设置 `Strict: true`。
11. TODO 11：原生模式不在 Prompt 中重复 JSON 示例，只保留问题语义和应用端校验。

## 9C：两种模式对比

12. TODO 12：切换 `selectedOutputMode`，分别运行 Prompt 与原生模式。
13. TODO 13：离线测试未知字段、多个 JSON、缺失字段、空数组和范围错误。
14. TODO 14：记录模型服务不支持 JSON Schema 时的错误，不静默降级到 Prompt 模式。

## 开始练习

先在 `main.go` 顶部修改 OpenAI-compatible 模型常量：

```go
const (
	baseURL   = "http://localhost:8084/v1"
	apiKey    = "replace-with-your-api-key"
	modelName = "gpt-5.4-mini"
)
```

将示例中的占位值替换为本地实际配置，不要把真实 API Key 写进源码或提交到 Git。

通过 `selectedOutputMode` 选择本次运行方式；每次只调用一次模型，避免对比练习产生双倍请求：

```go
const selectedOutputMode structuredOutputMode = promptOnlyMode
```

或：

```go
const selectedOutputMode structuredOutputMode = nativeJSONSchemaMode
```

```bash
go run ./examples/ai/phase2/10_structured_json_output
```

默认 API Key 或外部服务配置保持占位符时，程序应明确提示配置未完成，不得连接外部服务或产生付费请求。

## 验证方式

```bash
gofmt -w examples/ai/phase2/10_structured_json_output/*.go
go test -timeout=60s ./examples/ai/phase2/10_structured_json_output
go vet ./examples/ai/phase2/10_structured_json_output
```

## 完成标准

- TODO 编号连续，错误使用 `%w` 保留底层原因。
- Context 超时和取消能够传播到模型或网络边界。
- 不打印 API Key、完整敏感 Prompt 或 Authorization。
- Prompt 模式不设置 `ResponseFormat`，原生模式设置严格 JSON Schema。
- 两种模式都必须经过相同的 Decoder 和业务字段校验。
- 删除核心实现后，测试或默认运行能够清楚暴露练习未完成。

## 暂不实现

- 自动修复和自动重试循环。
- 服务端不支持 JSON Schema 时的静默降级。
