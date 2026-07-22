# Phase 2 综合生产化练习

## Goal

新增阶段 2 的第 13 个综合练习，把前 12 个练习中值得重复掌握的 AI 能力串成一条贴近生产的完整业务流程。练习应让学习者主要编写模型、Prompt、会话、结构化输出、流式输出和调用治理逻辑；通用 HTTP、配置装配等模板化代码可以预先提供。

## Background

- 阶段 2 当前有 12 个练习，目录编号为 `01` 到 `12`，目录契约测试当前固定期望 12 个。
- 已有知识覆盖：OpenAI-compatible 基础调用、Eino Generate、消息角色、Zero-shot/Few-shot、多轮会话、Redis 历史、历史截断/TTL/用户隔离、ChatTemplate、Stream、JSON Schema 结构化输出、调用治理和模型 Provider 适配。
- 当前练习多为单知识点 Demo；新练习应提供接近真实应用的端到端数据流，同时避免让学习者把主要时间消耗在样板代码上。
- 项目是电商 Go 服务，综合练习放在 `examples/ai/phase2/13_*` 最符合现有学习路径；生产业务代码不应被练习实现侵入。
- 真实模型练习必须保留 Eino/OpenAI-compatible 构造入口、明确真实调用和费用风险，并支持离线测试核心逻辑。
- 项目已有 Gin 和 Redis Client；阶段 2 的流式练习当前只写标准输出，没有可复用的 HTTP/SSE 业务接口。
- 阶段 2 练习 12 仍是 Provider 适配骨架，因此综合练习需要提供一套可运行的统一 Generate/Stream Provider 边界，而不是假设已有生产实现。
- 当前固定依赖为 Eino `v0.9.12`、Eino OpenAI-compatible adapter `v0.1.13`、Gin `v1.10.1` 和 go-redis `v9.12.1`。现有练习已证明本地 OpenAI-compatible 服务可经 `http://localhost:8084/v1` 接入，但密钥必须来自环境变量，不能写入源码或文档。

## Requirements

- R1：新增阶段 2 第 13 个综合练习，并同步目录导航、路线图和目录契约测试。
- R2：练习必须形成单一真实业务流程，不得只是把 12 个独立 Demo 复制到一个目录。
- R3：流程必须复习基础配置、Eino 模型调用、System/User/Assistant 消息、Few-shot、会话历史、Redis 隔离与截断、ChatTemplate、流式输出、结构化输出、调用治理和 Provider 适配。
- R4：学习者 TODO 聚焦 AI 相关决策与实现；通用 HTTP 路由、请求解码、配置装配、SSE framing 等低学习价值样板可以提供完整实现。
- R5：真实完成路径使用本地 OpenAI-compatible 服务；默认或测试路径不得意外产生付费请求。
- R6：模型失败时不得保存半轮会话；用户与 session 必须隔离；历史必须有长度/TTL 边界。
- R7：流式响应必须处理 EOF、取消、中途错误和资源关闭；结构化输出必须经过严格解析和业务校验。
- R8：模型调用必须有超时、限流、有限重试、可取消退避、错误白名单和延迟/Token/成功指标。
- R9：至少提供离线测试覆盖关键成功与失败路径；真实模型烟测不能替代离线行为测试。
- R10：README 必须说明业务流程、知识点映射、哪些代码已提供、学习者 TODO 顺序、真实运行依赖、费用风险和验收命令。
- R11：综合练习以电商智能客服会话服务为唯一业务主线，不提供互不相干的平行 Demo。
- R12：综合练习以 HTTP + SSE 服务运行；Gin 路由、请求解码、基础校验和 SSE framing 由骨架完整提供，学习者主要实现 AI 数据流。
- R13：每条用户消息采用两阶段模型调用：先用 Few-shot + JSON Schema `Generate` 得到意图/风险/回答风格，再用分类结果、ChatTemplate 和 Redis 历史执行 `Stream` 返回最终回答。
- R14：README 必须明确两阶段调用会增加延迟与 Token 成本，指标需要区分分类调用和回答调用，并提供请求级汇总。
- R15：真实服务强依赖 Redis，启动时 Ping 失败必须明确退出；不得静默降级为进程内历史或第二事实源。
- R16：离线测试通过 Fake History Store 验证业务逻辑，不要求连接 Redis；Fake 不得作为真实完成证明。
- R17：流式中途错误时不得保存 User 或半截 Assistant；服务发送 SSE `error` 事件并结束请求。
- R18：Stream 已输出首个文本 chunk 后禁止自动重试；只有创建 Stream 失败且尚未向客户端输出内容时，才允许按错误白名单有限重试。
- R19：本地服务只监听 `127.0.0.1`，通过已实现的开发身份中间件校验 `X-Demo-User-ID` 并写入 Context；AI Service 不得从请求 Body 读取用户身份。
- R20：README 必须明确开发身份 Header 不能用于生产，真实系统应由 JWT、网关或其他可信认证层注入身份；本练习不实现认证协议。
- R21：在 `knowledge.go` 内置确定性业务知识，至少覆盖商品建议、配送、退换货和售后升级条件，不依赖外部知识文件。
- R22：使用确定性选择器根据结构化意图选择业务上下文；不在本练习实现数据库查询、向量检索或 Tool Calling。
- R23：交付引导式练习骨架：通用基础设施和模板代码完整提供，AI 核心保留连续、可验证的 TODO；未完成时编译通过但真实运行在发起外部请求前明确返回 `errExerciseIncomplete`。
- R24：练习 TODO 必须映射到前 12 个练习的核心能力，避免把同一知识点拆成无意义的机械填空；README 明确哪些文件无需学习者修改。

## Acceptance Criteria

- [ ] `examples/ai/phase2` 包含连续编号的第 13 个练习，仓库目录契约测试通过。
- [ ] 一个用户请求能够沿统一服务流程完成意图理解、历史装配、模型生成、输出返回和成功后历史提交。
- [ ] 12 个前置练习的核心能力在 README 中映射到具体代码边界和可验证 TODO，不存在纯文字“已覆盖”。
- [ ] 普通生成、结构化生成和流式生成至少以两种真实业务分支接入同一 Provider/治理边界。
- [ ] Redis 不可用、模型不可重试错误、可重试错误耗尽、Context 取消、流中途错误、结构化解析失败均有明确失败行为。
- [ ] 离线测试、`go test -timeout=60s ./examples/ai/...` 和 `go vet ./examples/ai/...` 通过。
- [ ] README 指定的真实运行命令能够使用本地 OpenAI-compatible 模型完成闭环，且不会记录 API Key、Authorization 或完整敏感 Prompt。

## Out Of Scope

- 阶段 3 的 Embedding、向量数据库、RAG 和引用生成。
- 阶段 4 的 Tool Calling、订单查询工具和写操作确认。
- 完整生产部署、监控平台、分布式限流和多实例会话一致性。
- 为练习改造现有生产 Handler、Service 或 Repository。
