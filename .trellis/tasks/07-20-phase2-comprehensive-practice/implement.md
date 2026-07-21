# Phase 2 综合生产化练习：实施计划

## 实施前置条件

- [ ] 用户审核并批准 `prd.md`、`design.md` 和本计划后，执行 `task.py start`；在此之前不创建练习源码。
- [ ] 进入实现前读取 `.trellis/spec/` 中 `examples/ai` / Go 后端相关规则，并检查当前工作树，避免覆盖用户的 Embedding 规划任务。
- [ ] 确认现有 `examples/ai/phase2/07`、`08`、`09`、`10`、`11`、`12` 的接口与测试写法；使用现有 Eino `v0.9.12` API，不为示例升级依赖。

## 有序实施清单

1. 创建 `examples/ai/phase2/13_customer_support_service/` 和完整的 `config.go`、`identity.go`、`http.go`、`knowledge.go`、`main.go` 基础设施。
   - 提供 loopback 地址校验、开发身份 Context、请求校验、SSE writer、Redis Ping 和关闭路径。
   - 提供 `testdata/business_context.json`，并在启动时严格校验四个 intent 的知识数据。
   - Redis 历史只实现一个 Store；没有 in-memory fallback。

2. 定义与 SDK 隔离的模型、流、调用、Usage、指标和错误类型。
   - 在 `provider.go` 放置连续 TODO 1–2，提供 OpenAI-compatible Eino adapter 的清晰入口。
   - 以 `errExerciseIncomplete` 阻断未完成 Provider，在任何外部模型调用前返回。

3. 添加分类阶段骨架。
   - 在 `classification.go` 放置 TODO 3、4、10、11：角色化 Few-shot 消息、原生 JSON Schema、严格解码及 enum / 布尔字段校验。
   - 创建分类专用 Provider 实例；模型不支持 JSON Schema 时显式失败，不切换到另一个输出模式。

4. 添加回答 Prompt 与会话状态骨架。
   - 在 `history.go` 放置 TODO 6–7，使用隔离 Key、成对消息、`TxPipelined`、截断和 TTL。
   - 在 `service.go` 放置 TODO 5、8、9、15，使用 ChatTemplate 装配 System / History / User，消费流并执行“EOF 后再提交”的不变量。
   - 首个 delta 前后分别建立明确的失败分支，避免重试或保存半轮历史。

5. 添加统一调用治理骨架。
   - 在 `governance.go` 放置 TODO 12–14：每次尝试限流、独立超时 Context、错误白名单、上限指数退避、流 Context 生命周期及阶段指标。
   - Generate 和 Stream 创建均使用同一治理边界；只允许回答流创建失败重试。

6. 实现并整理引导文档和目录契约。
   - README 映射 TODO 1–15 到阶段 2 的前 12 题，标明“已提供”和“必须完成”的文件。
   - 说明环境变量、Redis 强依赖、真实调用的 Token / 延迟成本、开发 Header 的边界、真实烟测和离线验证。
   - 更新 `examples/ai/README.md`、路线图的阶段 2 顺序和数量，并将 `catalog_test.go` 的 Phase 2 计数更新为 13；只为第 13 题增加跨文件 TODO / README 对齐校验。

7. 编写离线测试和真实烟测说明。
   - 使用 Fake Provider 和 Fake History Store，不启动 Redis、不访问模型。
   - 先覆盖未完成骨架的无外部调用行为；在 TODO 完成后自动启用成功和失败路径断言。
   - 通过 `httptest` 验证 Header 身份、SSE 事件顺序和失败不提交。

## 验证命令

实现期间先运行最小目标测试，再运行全量 AI 示例校验：

```bash
go test -timeout=60s ./examples/ai/phase2/13_customer_support_service
go test -timeout=60s ./examples/ai/...
go vet ./examples/ai/...
git diff --check
```

README 还会提供一个不含敏感信息的真实验证流程：启动 Redis 与本地 OpenAI-compatible 模型后运行服务，并用 `curl -N` 发送一条 SSE 请求。真实烟测需要显式配置 `OPENAI_API_KEY`；离线测试永不读取或打印该值。

## Review 与回滚

- 执行后使用 `trellis-check` 和 `go-review` 进行主审查；因本题涉及 Redis、一致性、身份边界和流式取消，按项目规则启动只读独立审查并复验修复。
- 重点检查：Stream 成功后是否过早取消 Context、首个 delta 后是否仍会重试、历史是否可能留下半轮、SSE `done` 是否早于 Redis 提交、日志是否泄漏 Prompt 或凭证、目录契约是否削弱其他练习校验。
- 如果实现阶段发现阶段 2 的旧练习接口无法支撑设计，只回到本任务的设计文档调整；不改造现有生产业务代码。
- 回滚范围仅为第 13 题目录、三处导航 / 计数变更和本任务文档；不触及主服务模块或数据库迁移。
