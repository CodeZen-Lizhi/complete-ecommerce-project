# 阶段 2 练习 13：电商智能客服综合服务

## 练习目标

把阶段 2 的模型调用、消息角色、Few-shot、会话历史、ChatTemplate、流式输出、结构化输出、调用治理和 Provider 适配串成一个独立的电商智能客服 HTTP + SSE 服务。

服务只监听回环地址，唯一接口为：

```text
POST /api/ai/chat/stream
X-Demo-User-ID: <开发环境身份>
Content-Type: application/json

{"session_id":"checkout-help","message":"这款耳机多久能到？"}
```

身份只能来自已验证的 `X-Demo-User-ID` Header；请求体中的 `user_id` 会被严格拒绝。真实生产系统应由 JWT、网关或其他可信认证层注入身份，本练习不实现认证协议。

## 服务流程

1. 在开始 SSE 前校验 Header、`session_id` 和 `message`。
2. 用 Few-shot + 严格 JSON Schema 分类意图、回答风格和是否需要人工处理。
3. 根据已校验意图选择本地业务知识，读取当前用户当前会话的 Redis 历史。
4. 使用 ChatTemplate 组装回答消息并启动回答流。
5. 仅在流正常 EOF、回答非空且 Redis 原子保存完整 User/Assistant 轮次后发送 `done`。

流内事件只有 `meta`、`delta`、`done`、`error`。`meta` 包含 `intent`、`response_style`、`requires_handoff`；`done` 包含分类与回答两个阶段各自的 attempts、延迟和 Token，以及请求级 `total_latency_ms`、`total_usage` 汇总。模型、Redis 或解析内部错误不得直接返回给客户端。

## 已提供的基础设施

- `config.go`：环境变量、回环地址和数值边界校验。
- `identity.go`、`http.go`：开发身份中间件、严格 JSON 解码、SSE framing 和安全错误映射。
- `knowledge.go`、`testdata/business_context.json`：固定意图和确定性业务知识选择。
- `main.go`：Redis Ping、两个独立 Provider 的装配、回环监听与优雅关闭。
- `main_test.go`：不连接 Redis 或模型的离线保护测试。

## TODO 顺序

1. TODO 1：在 `provider.go` 的 `Generate` 中调用真实 Eino 模型并归一化文本与 Token Usage。
2. TODO 2：在 `provider.go` 的 `Stream` 中封装 Eino Reader 的 `Recv`/`Close` 生命周期。
3. TODO 3：在 `classification.go` 写出严格分类 System Message。
4. TODO 4：在 `classification.go` 添加三个 User/Assistant Few-shot 分类样例。
5. TODO 5：在 `service.go` 只读取当前用户当前会话，并在最终成功后提交完整轮次。
6. TODO 6：在 `history.go` 构造隔离 Redis Key，读取末尾完整历史轮次。
7. TODO 7：在 `history.go` 使用 `TxPipelined` 原子追加、裁剪并刷新 TTL。
8. TODO 8：在 `service.go` 用 `ChatTemplate` 组装 System、History、User 消息。
9. TODO 9：在 `service.go` 消费流、处理 EOF/取消/错误，并只转发非空 delta。
10. TODO 10：在 `classification.go` 构建严格原生 JSON Schema ResponseFormat。
11. TODO 11：在 `classification.go` 用严格 Decoder 解析并校验分类结果。
12. TODO 12：在 `governance.go` 实现限流、尝试级超时和非流式指标。
13. TODO 13：在 `governance.go` 实现流创建前的有限重试和流 Context 生命周期。
14. TODO 14：在 `governance.go` 实现可重试错误白名单与可取消、有上限的退避。
15. TODO 15：在 `service.go` 核对两阶段编排、SSE 事件顺序和“成功后再提交历史”不变量。

未完成时，核心路径必须返回同一个 `errExerciseIncomplete`，并且在任何模型网络调用前停止。不得用 Fake Provider、内存历史或伪造回答掩盖未完成实现。

## 前置练习映射

| 本题 TODO | 对应前置练习 | 在本题中复习的能力 |
| --- | --- | --- |
| 1–2 | 01、02、12 | OpenAI-compatible 配置、Eino Generate/Stream 与 Provider 适配边界 |
| 3 | 03 | System、User、Assistant 消息角色 |
| 4 | 04 | 固定 Few-shot 分类样例 |
| 5 | 05 | 多轮会话的读取与最终成功后提交 |
| 6–7 | 06、07 | Redis 会话隔离、截断、TTL 与原子完整轮次 |
| 8 | 08 | ChatTemplate 与 typed 历史占位符 |
| 9 | 09 | Stream EOF、取消、中途错误与关闭 |
| 10–11 | 10 | 原生 JSON Schema、严格解码和业务枚举校验 |
| 12–14 | 11 | 限流、尝试级超时、重试白名单和可取消退避 |
| 15 | 12 | 两阶段 Provider 编排、SSE 事件顺序和提交不变量 |

## 运行前配置

| 环境变量 | 默认值 / 要求 |
| --- | --- |
| `AI_DEMO_LISTEN_ADDR` | `127.0.0.1:8093`，必须是数值回环地址 |
| `AI_DEMO_MODEL_BASE_URL` | `http://localhost:8084/v1` |
| `OPENAI_API_KEY` | 必填，仅从环境读取 |
| `AI_DEMO_MODEL` | `gpt-5.4-mini` |
| `AI_DEMO_REDIS_ADDR` | `127.0.0.1:6379` |
| `AI_DEMO_HISTORY_TTL` | `30m` |
| `AI_DEMO_HISTORY_TURNS` | `6` |
| `AI_DEMO_CALL_TIMEOUT` | `30s` |
| `AI_DEMO_MAX_RETRIES` | `2` |
| `AI_DEMO_INITIAL_BACKOFF` | `200ms` |
| `AI_DEMO_CALLS_PER_SECOND` | `2` |

完成 TODO 后，先启动本地 Redis 和支持 OpenAI-compatible JSON Schema、流式 Usage 的模型服务，再执行：

```bash
cd examples/ai/phase2/13_customer_support_service
go run .
```

另一个终端可发送：

```bash
curl -N \
  -H 'Content-Type: application/json' \
  -H 'X-Demo-User-ID: learner-1001' \
  -d '{"session_id":"checkout-help","message":"这款耳机多久能到？"}' \
  http://127.0.0.1:8093/api/ai/chat/stream
```

该流程包含一次分类调用和一次回答调用。可重试故障会增加延迟、Token 和模型成本；不得把外部、共享或生产 API Key 写入源码、README 或日志。

## 验证方式

```bash
gofmt -w examples/ai/phase2/13_customer_support_service/*.go
go test -timeout=60s ./examples/ai/phase2/13_customer_support_service
go test -timeout=60s ./examples/ai/...
go vet ./examples/ai/...
```

## 完成标准

- TODO 1–15 在五个核心文件中唯一、连续，并与本 README 对齐。
- 分类失败、流创建失败、流中途错误、客户端取消、空回答和历史提交失败时均不写入半轮历史。
- 首个 `delta` 后绝不自动重启流；只有尚未输出内容的 Stream 创建失败可按白名单重试。
- 日志不包含 API Key、Authorization、完整 Prompt、完整用户消息或模型原始输出。
- Redis 不可用时进程启动失败；不得降级为进程内历史。

## 暂不实现

- 订单数据库查询、向量检索、RAG、Tool Calling 和写操作。
- 生产认证、分布式限流、多实例会话一致性和监控平台。
