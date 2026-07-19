# 阶段 2 练习 4：Zero-shot 与 Few-shot 对比

## 练习目标

使用真实 Eino `AgenticModel` 对同一个分类问题分别进行 Zero-shot 和 Few-shot 调用，观察示例消息怎样影响输出格式、标签选择和稳定性。

本练习只改变一个变量：是否在最终 User Message 前提供成对的 User/Assistant 示例。模型、System Prompt、最终问题和调用参数必须完全相同。

## 知识点

- Zero-shot、One-shot 与 Few-shot 的区别。
- Few-shot 示例必须使用真实消息角色，而不是把所有示例拼成一段说明文字。
- 示例顺序、标签覆盖范围和错误示例会影响模型泛化。
- Few-shot 能提高格式遵循率，但不能替代 JSON Schema 或应用端校验。

## TODO 顺序

1. TODO 1：准备一个分类问题和两组 User/Assistant 示例。
2. TODO 2：构造 Zero-shot 与 Few-shot 消息，保证最终问题一致。
3. TODO 3：从顶部常量读取配置并创建真实 Eino ResponsesModel。
4. TODO 4：使用相同模型依次执行两个场景。
5. TODO 5：严格提取非空 AssistantGenText。
6. TODO 6：比较格式遵循率、标签稳定性、照抄和错误泛化。

## 消息结构

```text
Zero-shot: System -> 最终 User

Few-shot:  System
           -> 示例 User 1 -> 示例 Assistant 1
           -> 示例 User 2 -> 示例 Assistant 2
           -> 最终 User
```

## 开始练习

先把 `main.go` 顶部 `apiKey` 占位值替换为本地实际配置，再运行：

```bash
go run ./examples/ai/phase2/04_few_shot_comparison
```

该命令会产生真实模型调用。真实 Key 只在本地临时使用，不要提交到 Git。

## 观察重点

- Few-shot 是否更稳定地只返回允许标签。
- 模型是否机械照抄示例中的解释或字段值。
- 示例没有覆盖的新输入能否正确泛化。
- Context 取消和模型错误是否被明确传播。

## 完成标准

- 两组调用除示例消息外完全一致。
- Few-shot 使用独立 User/Assistant 消息对，不是单段字符串拼接。
- 输出真实模型回答并能解释示例对结果的正负影响。
- 占位 Key、空响应和 Context 超时都明确失败，不返回假结果。

## 暂不实现

- 不用 Fake/Mock 结果作为完成证明。
- 不在本练习实现 JSON Schema；原生结构化输出由后续练习负责。
- 不自动搜索或生成 Few-shot 示例。

## 验证

```bash
gofmt -w examples/ai/phase2/04_few_shot_comparison/main.go
go test -timeout=60s ./...
go vet ./...
```
