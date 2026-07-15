# 阶段 2：Redis 会话历史

## Goal

在上一练习的 Eino ResponsesModel 多轮对话基础上，将消息历史从进程内切片迁移到 Redis。每次处理用户输入时，都按 session ID 重新读取并组装完整 Messages，使第二轮调用不依赖当前进程内存也能理解上一轮上下文。

## Background

- 学习路线阶段 2 的第 5 项是“使用 Redis 按 session 保存历史，下一轮重新组装 Messages”。
- 上一练习已掌握 `System -> User 1 -> Assistant 1 -> User 2` 的内存历史组装方式，本题只替换历史存储边界。
- 仓库已依赖 `github.com/redis/go-redis/v9 v9.12.1`，本地开发编排已提供 Redis 7，不新增依赖。
- Eino 继续使用 `agenticopenai.ResponsesModel` 与 `schema.AgenticMessage`，并关闭服务端存储和自动缓存，确保上下文来自 Redis。

## Requirements

- 新建独立示例目录 `examples/ai/phase2/redis_session_history/`，可通过 `go run ./examples/ai/phase2/redis_session_history` 运行。
- 使用 Redis List 保存会话历史：一个经过校验的 session ID 对应一个带固定命名空间的 Redis Key，每个 List 元素保存一条 JSON 编码的完整 `schema.AgenticMessage`。
- 读取历史时使用 `LRANGE 0 -1`，按 Redis List 原顺序逐条反序列化；空 List 视为新会话，不视为错误。
- 追加历史时先校验输入并逐条 JSON 序列化，再通过一次 `RPUSH` 按原顺序追加，避免消息顺序被打乱。
- 将单轮处理抽象为“读取历史 -> 新会话补 System -> 追加当前 User -> 调用模型 -> 保存本轮新增消息 -> 输出 Assistant”；主程序用同一 session ID 调用两次，第二轮 Prompt 必须依赖第一轮内容。
- 只持久化本轮新增的 System（仅新会话）、User 和 Assistant，不能把已读取的全部历史再次追加造成重复。
- Redis、模型调用和 JSON 编解码错误必须带操作与 session 上下文，并使用 `%w` 保留底层错误链；模型与 Redis 调用继续接收并传播 Context。
- API Key 只保留 `replace-with-your-api-key` 占位符；占位符未替换时明确失败，不连接 Redis、不调用模型。
- 骨架保留必要常量、接口、类型、模型/Redis Client 创建和非核心辅助方法；Redis 历史读写及单轮组装由学习者按编号连续的中文 TODO 完成。
- 未完成 TODO 时返回明确的练习未完成错误，不产生真实 Redis 写入或付费模型调用。

## Out Of Scope

- 不限制历史轮数或 Token，不做摘要、裁剪和上下文预算。
- 不设置或刷新会话 TTL，不处理会话过期。
- 不实现 user/tenant 与 session 的绑定、跨用户授权或同一 session 的并发写冲突。
- 不使用 Stream、ChatTemplate、结构化输出、重试、限流或调用统计。
- 不接入电商 HTTP 路由、Container 或业务 Service，不修改现有 Redis 单例封装和已有练习。

## Acceptance Criteria

- [x] 骨架经过 `gofmt` 且目标包可编译；API Key 为占位符或 TODO 未完成时，不连接 Redis、不调用真实模型。
- [x] 中文 TODO 编号连续，完整覆盖 session Key、`LRANGE`、JSON 反序列化、JSON 序列化、`RPUSH`、新会话 System 初始化、本轮消息保存和两次单轮调用，但不预先给出核心答案。
- [x] 学习者完成 TODO 后，同一 session 的第二次调用会从 Redis 读取 `System -> User 1 -> Assistant 1`，追加 User 2 后再调用模型。
- [x] Redis 中每条历史消息只出现一次，角色和消息顺序与模型实际调用顺序一致。
- [x] 空历史可初始化新会话；空 session ID、nil 消息、损坏 JSON、Redis 错误、空模型响应均明确失败。
- [x] 源码、任务文件和提交内容不包含真实 API Key。
- [x] `go test -timeout=60s ./examples/ai/phase2/redis_session_history`、`go vet ./examples/ai/phase2/redis_session_history` 和 `git diff --check` 通过。

## Technical Notes

- Redis List 天然保留追加顺序，适合本题观察消息序列；本题不把它扩展为通用生产会话存储方案。
- `schema.AgenticMessage` 及其内容块提供 JSON 标签，可以直接编码完整消息；展示文本仍由 `assistantText` 单独提取。
- 当前任务是单目录学习骨架，采用 PRD-only 轻量规划，不新增 `design.md` 或 `implement.md`。
