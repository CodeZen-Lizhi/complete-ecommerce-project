# 阶段 2：程序内多轮对话

## Goal

在上一练习的 Eino ResponsesModel 基础上，用进程内消息切片完成两轮连续对话，理解模型 API 默认无状态、应用必须保存 Assistant 回答并在下一次调用时重新发送完整历史。

## Background

- 学习路线阶段 2 的第 4 项是“使用消息列表完成程序内多轮对话”。
- 上一练习已经使用 `einomodel.AgenticModel`、`schema.AgenticMessage` 和 `agenticopenai.NewResponsesModel` 完成 System/User Prompt 对照。
- Eino ResponsesModel 支持 `EnableAutoCache`，但本练习必须关闭该能力，避免服务端缓存掩盖应用侧历史重组过程。
- 本练习是独立学习示例，不接入电商业务入口、Container、Redis 或数据库。

## Requirements

- 新目录为 `examples/ai/phase2/in_memory_multi_turn/`，只保留 `main.go`，可以通过 `go run ./examples/ai/phase2/in_memory_multi_turn` 独立运行。
- 继续使用 Eino `agenticopenai.ResponsesModel` 和 `schema.AgenticMessage`，不切换回 Chat Completions。
- 使用同一个模型实例完成两次顺序调用，消息历史保存在当前进程的 `[]*schema.AgenticMessage` 中。
- 初始历史只包含一条 System Message；第一轮追加 User Message，调用模型后再追加返回的 Assistant Message。
- 第二轮 User Prompt 必须依赖第一轮上下文，追加后将完整历史再次传给 `Generate`，以证明模型回答来自应用重组的上下文。
- `ResponsesConfig` 明确设置 `Store=false`、`EnableAutoCache=false` 和 `MaxRetries=0`；不使用 response ID 或服务端对话状态。
- 每次模型调用使用独立的 30 秒子 Context，并及时调用 `cancel`，避免第一轮占用第二轮的超时预算。
- API Key 在源码中只保留 `replace-with-your-api-key` 占位符；占位符未替换时程序明确失败且不发起模型请求。
- 骨架保留必要常量、类型、函数签名、模型创建和非核心辅助代码；核心历史组装与两轮调用由学习者按编号连续的中文 TODO 完成。
- TODO 必须按执行顺序覆盖消息初始化、追加 User、调用 Generate、追加 Assistant、追加追问、再次调用、输出结果和错误处理。

## Out Of Scope

- 不接入 Redis，不实现 session ID、并发用户隔离、会话过期或持久化。
- 不做历史轮数/Token 截断；这些属于阶段 2 第 6 个练习。
- 不使用 `EnableAutoCache`、`previous_response_id` 或服务端保存的 Response 历史。
- 不实现 Stream、ChatTemplate、结构化 JSON、重试、限流和 Token 统计。
- 不创建 README 或测试文件；学习说明直接保留在中文 TODO 和阶段路线中。
- 不修改已有 `basic_chat`、`eino_generate`、`prompt_roles` 或电商业务代码。

## Acceptance Criteria

- [x] `main.go` 经过 `gofmt` 并可以编译；Key 仍为占位符时不调用真实模型。
- [x] 中文 TODO 编号连续、步骤完整，但没有预先给出核心多轮对话答案。
- [x] 学习者完成 TODO 后，程序按 `System → User 1 → Assistant 1 → User 2` 保存历史，并用完整列表生成 Assistant 2。
- [x] 第二轮追问依赖第一轮内容，输出分别标识两轮 User Prompt 和 Assistant 回答，能够观察上下文是否生效。
- [x] 两次调用分别受 30 秒超时控制，错误包含轮次上下文并保留底层错误链。
- [x] 提交前源码和任务文件中不包含真实 API Key。
- [x] `go build ./examples/ai/phase2/in_memory_multi_turn` 与 `go vet ./examples/ai/phase2/in_memory_multi_turn` 通过。

## Technical Notes

- 完整 Assistant `AgenticMessage` 包含角色、文本块及可能的扩展信息，第一轮响应应作为历史消息整体追加，不应只保存打印出来的字符串。
- `assistantText` 只负责展示时提取 `AssistantGenText`；历史保存与展示文本是两个不同职责。
- 本任务局限于一个独立示例，采用 PRD-only 轻量规划，不新增 `design.md` 或 `implement.md`。
