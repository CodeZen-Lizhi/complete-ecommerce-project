# 阶段 2：会话历史截断与隔离

## Goal

在上一练习的 Redis 会话历史基础上，生成一个可独立运行的阶段 2 第 6 个学习骨架，让学习者亲手实现“只发送最近若干轮历史、会话自动过期、不同用户之间严格隔离”。练习继续强调模型 API 无状态，应用负责保存、裁剪和重组上下文。

## Background

- 学习路线 `docs/ai-learning/go-eino-ai-engineer-roadmap.md` 的阶段 2 第 6 项要求“限制历史轮数或 Token，处理会话过期和用户隔离”。
- 上一练习 `examples/ai/phase2/05_redis_session_history/main.go` 已完成 Redis List 持久化，但 Key 只包含 session ID、读取全部历史且没有 TTL。
- 本练习采用“限制历史轮数”而不是精确 Token 计数，避免提前引入 tokenizer 或模型厂商差异；Token 预算留作后续深化。
- Redis 只保存 User/Assistant 消息对；System Prompt 每次调用时由应用重新注入。这样可以安全使用 `LTRIM` 保留末尾若干条消息，而不会误删唯一的 System Message。

## Requirements

- 新建独立目录 `examples/ai/phase2/06_session_history_limits/`，不接入电商 HTTP 服务、Container 或生产 Redis 封装。
- 提供可编译、可运行的 Go 骨架和中文 README；核心学习步骤使用执行顺序连续的中文 TODO，不预先给出完整答案。
- 使用 Eino `agenticopenai.ResponsesModel`，关闭服务端 Store 和自动缓存，API Key 只保留占位符。
- Redis Key 必须同时包含经过校验的 user ID 和 session ID，并使用固定命名空间，防止同名 session 跨用户读取。
- Redis List 只保存 User/Assistant 消息；每轮成功后原子执行追加、按消息上限裁剪和刷新 TTL，模型失败时不得保存半轮对话。
- 每次模型调用前从 Redis 重新加载最近历史，并按 `System -> 最近历史 -> 当前 User` 组装 Messages。
- 历史上限以完整轮次表示，每轮固定两条消息；配置必须拒绝非正轮数和非正 TTL。
- 未完成核心 TODO 时返回明确的练习未完成错误；默认占位 API Key 下不得连接 Redis 或发起真实模型调用。
- 所有新增 Go 函数和方法写简洁中文注释，错误使用 `%w` 保留底层原因，不打印 API Key。

## Out of Scope

- 精确 Token 计数、摘要压缩、语义记忆或向量检索。
- Redis Lua 脚本、分布式锁、并发写冲突治理和生产级多租户鉴权。
- ChatTemplate、Stream、结构化 JSON、重试、限流与调用统计；这些属于后续阶段 2 练习。
- 修改上一练习或生产业务代码。

## Acceptance Criteria

- [x] `go test -timeout=60s ./examples/ai/phase2/06_session_history_limits` 与 `go vet ./examples/ai/phase2/06_session_history_limits` 通过。
- [x] 默认配置运行时明确提示 API Key 未配置，不连接 Redis、不产生付费模型请求。
- [x] 骨架中的 TODO 编号连续，覆盖身份校验与 Key、读取最近历史、消息重组、原子追加/裁剪/TTL、两轮隔离演示和错误处理。
- [x] README 解释为什么 Redis 不保存 System Message、为什么按偶数消息裁剪，以及如何验证用户隔离和 TTL。
- [x] 练习只新增目标目录和本任务必要的 Trellis 文件，不改动现有业务与前序练习。
