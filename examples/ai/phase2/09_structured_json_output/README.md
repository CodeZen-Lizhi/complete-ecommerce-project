# 阶段 2 练习 9：结构化 JSON 输出

## 练习目标

通过 Prompt 约束模型只返回固定 JSON，并完成严格解析和业务字段校验。

## TODO 顺序

1. TODO 1：编写禁止 Markdown 代码块和额外文本的 System Prompt。
2. TODO 2：在 User Message 中提供固定 JSON 结构和当前问题。
3. TODO 3：通过 Model Client 生成原始文本，模型错误不得当成空 JSON。
4. TODO 4：把原始文本直接交给严格解析器，不静默截取或修复。
5. TODO 5：使用 Decoder 严格解析并拒绝未知字段。
6. TODO 6：拒绝第二个 JSON 值和尾随非空内容。
7. TODO 7：校验必填、数组元素和 Confidence 范围。
8. TODO 8：区分语法、结构和业务字段错误，并保留错误链。

## 开始练习

```bash
go run ./examples/ai/phase2/09_structured_json_output
```

默认 API Key 或外部服务配置保持占位符时，程序应明确提示配置未完成，不得连接外部服务或产生付费请求。

## 验证方式

```bash
gofmt -w examples/ai/phase2/09_structured_json_output/*.go
go test -timeout=60s ./examples/ai/phase2/09_structured_json_output
go vet ./examples/ai/phase2/09_structured_json_output
```

## 完成标准

- TODO 编号连续，错误使用 `%w` 保留底层原因。
- Context 超时和取消能够传播到模型或网络边界。
- 不打印 API Key、完整敏感 Prompt 或 Authorization。
- 删除核心实现后，测试或默认运行能够清楚暴露练习未完成。

## 暂不实现

- JSON Schema 原生响应格式和自动修复循环。
