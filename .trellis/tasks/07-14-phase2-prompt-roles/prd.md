# 阶段 2：System 角色与 User Prompt 实验

## Goal

在已完成基础 HTTP 调用和 Eino `Generate` 调用后，通过对同一个 User Prompt 使用两种不同的 System Prompt，观察角色、受众和输出约束如何改变模型回答。该练习对应学习路线阶段 2 的第 3 项，不进入多轮对话、Redis、Stream 或 RAG。

## Background

- `examples/ai/phase2/basic_chat` 已覆盖基础模型调用。
- `examples/ai/phase2/eino_generate` 已覆盖 Eino ChatModel `Generate`。
- 现有 `eino_generate` 虽然包含 System/User Message，但没有形成可对照的角色实验。
- 本练习继续采用“可编译骨架 + 编号连续的中文 TODO”，核心实现由学习者完成。

## Requirements

- 新练习目录为 `examples/ai/phase2/prompt_roles/`，不创建 `cmd/` 或阶段 3 目录。
- 使用 Eino 官方 `agenticopenai.NewResponsesModel` 请求 OpenAI-compatible Responses API。
- 对同一个 User Prompt 构造两组消息：一组面向 Go 初学者，一组面向有经验的 Go 工程师。
- 两组消息只改变 System Prompt，User Prompt 必须完全相同，保证对照变量单一。
- 完成后按顺序调用 `Generate` 两次，并分别输出场景名称和模型回答。
- 调用使用同一个带 30 秒超时的 `context.Context`，错误必须带场景上下文并保留错误链。
- 本地 AI 网关的自定义 API Key 允许直接配置在学习示例中，不在输出和错误中打印该值。
- 初始骨架在 TODO 未完成时明确失败，不发起真实模型调用。
- README 解释 System Prompt、User Prompt、控制变量和建议观察点，不提前给出 TODO 的完整代码答案。

## Out of Scope

- 不保存 Assistant 历史，不实现程序内多轮对话。
- 不接入 Redis，不处理 session、历史截断或过期。
- 不实现 ChatTemplate、Stream、结构化 JSON、重试、限流或调用统计。
- 不修改 `basic_chat`、`eino_generate` 或电商业务代码。

## Acceptance Criteria

- [ ] `go run ./examples/ai/phase2/prompt_roles` 可以独立运行。
- [ ] 初始骨架可编译，TODO 未完成时返回明确错误且不调用模型。
- [ ] 完成 TODO 后，同一个 User Prompt 会在两种 System Prompt 下分别生成回答。
- [ ] 输出明确标识“初学者导师”和“资深工程师顾问”两个场景，便于比较。
- [ ] Responses 模型创建、Generate 调用、HTTP 失败和空输出均有明确错误处理。
- [ ] 仅包含用户明确允许的本地 AI 网关 Key，不包含其他真实凭据；只新增 Eino 官方 `agenticopenai` 组件依赖。
- [ ] `go test -timeout=60s ./...` 与 `go vet ./...` 通过。

## Technical Notes

- 使用 `einomodel.AgenticModel`、`schema.AgenticMessage` 和 Eino 官方 `agenticopenai.NewResponsesModel`。
- Responses 模型文本从 `AgenticMessage.ContentBlocks` 中的 `AssistantGenText` 提取。
- 练习的核心变量只有 System Prompt；模型、User Prompt、Temperature 和调用方式保持一致。
- 两次调用按顺序执行，第一版不引入并发，以便初学者观察完整控制流和错误路径。
