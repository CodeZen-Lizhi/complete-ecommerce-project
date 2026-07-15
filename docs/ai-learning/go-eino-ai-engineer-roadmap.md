# Go + Eino AI 应用工程师学习路线（企业实用版）

> 目标：从 Go 初学和小 Demo 起步，最终达到能够参加 AI 应用工程岗位面试，并能在企业项目中开发、排查和优化 AI 系统的水平。

## 使用方法

每个阶段只保留四项：

1. **阶段目标**：这一阶段解决什么问题。
2. **必学理论**：需要查资料理解的知识。
3. **练习顺序**：必须亲手编码掌握的能力，可据此自行完成小 Demo。
4. **完成标志**：判断能否进入下一阶段。

通用原理和企业常用方案必须学习；少见、学术性强或岗位相关的方案放在文末“进阶专题”，知道用途即可。

### 练习协作规范

从阶段 2 当前练习开始，后续所有编码练习默认采用“骨架 + 中文 TODO”方式：

1. 保留完成练习所需的结构体、接口、方法签名、入口和必要的非核心代码。
2. 核心实现由学习者亲手完成，不预先提供完整答案；骨架应保持可编译、可运行，并在未完成时给出明确提示。
3. 在待实现位置按真实执行顺序编写编号连续的中文 TODO，步骤必须完整说明输入、处理、输出、错误处理、资源释放和关键边界校验。
4. 学习者按 TODO 逐步实现；默认只对当前步骤进行概念讲解、代码检查和错误诊断，不直接补齐后续实现。
5. 学习者完成后进行 review，指出具体问题和原因；明确问题由学习者优先修改，只有明确要求时才直接给出修复代码。
6. API Key 等敏感信息只允许在本地临时使用，不写入学习文档、不提交到 Git；示例和骨架只保留占位符。

### 编号练习目录

全部 7 个阶段、60 个练习已按阶段内顺序生成到 [examples/ai/README.md](../../examples/ai/README.md)。目录统一使用 `01_<slug>`、`02_<slug>` 格式；阶段 1 目录负责导航到现有电商代码，阶段 2–7 提供独立的 Go 学习骨架。

```text
Go 后端基础
→ LLM API 与 Prompt
→ RAG
→ Tool Calling
→ Agent
→ Eval
→ 安全、监控与部署
```

---

# 阶段 1：Go 后端基础

## 阶段目标

具备开发 AI 后端所需的 Go、HTTP、数据库、缓存、并发和测试能力。现有电商项目作为本阶段练习环境。

## 必学理论

- struct、interface、指针、slice、map 和错误处理。
- Context 的超时、取消和调用链传播。
- goroutine、channel、锁、数据竞争和 goroutine 泄漏。
- HTTP、JSON、REST、幂等和状态码。
- Handler、Service、Repository 分层。
- 数据库索引、事务、隔离级别、分页和 N+1。
- Redis Cache Aside、过期时间、穿透、击穿和雪崩。
- JWT 认证与业务授权的区别。
- 日志、测试、限流和优雅停机。

## 练习顺序

1. 使用 Gin 编写带参数校验和统一响应的接口。
2. 使用 GORM 完成 CRUD、分页、事务和慢查询 Explain。
3. 使用 Redis 缓存查询结果，处理空值和缓存失效。
4. 实现 JWT 登录、中间件鉴权和后端权限判断。
5. 使用 HTTP Client 调用外部服务，处理超时和非 2xx。
6. 使用 Worker Pool 控制并发、取消任务并收集错误。
7. 为 Service、Handler、数据库和失败路径编写测试。

## 完成标志

- 能独立修改当前 Gin/GORM 项目。
- 能解释 Context、事务、缓存一致性和接口幂等。
- 能处理并发、超时、错误和测试。

---

# 阶段 2：LLM API 与 Prompt

## 阶段目标

稳定调用模型，理解 Messages、对话历史、Prompt、Stream、结构化输出和模型调用治理。

## 必学理论

- System、User、Assistant Messages。
- Token、上下文窗口、Temperature 和输出不确定性。
- 模型 API 通常无状态，对话记忆由应用保存和重组。
- System Prompt、Few-shot、上下文和输出约束。
- 普通输出与 Stream 流式输出。
- API Key、Base URL、模型名、超时、限流和错误分类。
- 哪些错误可以重试，指数退避和最大重试次数。
- Token、首 Token 延迟、总耗时和成本。

## 练习顺序

1. 配置 Base URL、API Key、模型名，完成一次基础模型调用。
2. 使用 Eino ChatModel `Generate` 完成相同调用。
3. 加入 System 角色和 User Prompt，观察角色变化。
4. 使用消息列表完成程序内多轮对话。
5. 使用 Redis 按 session 保存历史，下一轮重新组装 Messages。
6. 限制历史轮数或 Token，处理会话过期和用户隔离。
7. 使用 Eino ChatTemplate 管理 Prompt 和变量。
8. 使用 ChatModel `Stream`，处理 EOF、取消和中途错误。
9. 要求模型返回固定 JSON，并完成解析和字段校验。
10. 增加 Context 超时、有限重试、限流和调用统计。
11. 使用项目自定义接口隔离具体模型厂商，支持配置切换。

## 完成标志

- 能实现普通调用、流式调用和 Redis 多轮对话。
- 能解释模型无状态、上下文截断和重试边界。
- 能使用 ChatTemplate、结构化输出和模型适配层。
- 能记录并控制延迟、Token 和成本。

---

# 阶段 3：RAG

## 阶段目标

掌握企业最常用的文档处理、切块、索引、混合检索、重排、引用、索引维护和评估能力。

## 必学理论

### 基础流程

```text
文档解析 → 切块 → Embedding → 索引
用户问题 → Query 处理 → 检索 → 重排 → 上下文 → 回答
```

- Loader、Transformer、Embedding、Indexer、Retriever 的职责。
- Embedding、向量维度、Cosine/Inner Product/L2 和 Top-K。
- RAG 的检索错误、生成幻觉和引用来源。

### 企业常用切块策略

1. **递归/Token 长度切块**：通用基线，控制大小和 overlap。
2. **结构化切块**：按标题、段落和列表切分。
3. **Parent-Child 父子切块**：小块精确召回，父块提供完整上下文。
4. **特殊内容处理**：表格、代码块和 FAQ 单独处理。

### 企业常用检索能力

- 普通 Chunk 向量索引。
- Parent-Child Retrieval。
- Dense Vector + BM25 的 Hybrid Search。
- Metadata Filter：租户、权限、时间和文档类型。
- HNSW：理解 `M`、`efConstruction`、`efSearch` 对召回、内存和延迟的影响。
- Query Rewrite、Multi-Query 和多轮问题改写。
- RRF 融合、Rerank、去重和 Token Budget。

### 数据与索引生命周期

- MySQL 保存文档、版本、状态和索引 ID；向量库保存 Chunk、向量和检索 Metadata。
- 幂等导入、增量更新、删除、失败补偿和最终一致性。
- Embedding 模型或维度变化后需要重建和切换索引。
- RAG 质量必须通过检索指标和回答指标验证。

## 练习顺序

1. 调用 Embedding，使用 Go 手工计算余弦相似度和 Top-K。
2. 加载 Markdown/PDF，保存标题、页码、来源和权限 Metadata。
3. 对比递归、结构化、Parent-Child 和表格特殊处理。
4. 使用 Eino Indexer 写入一种向量数据库，Retriever 完成检索。
5. 调整 HNSW 搜索参数，记录命中率和 P95 延迟。
6. 实现 Parent-Child：子块检索，命中后返回父块。
7. 实现 Dense + BM25，并使用 RRF 融合结果。
8. 实现多轮问题改写、Multi-Query 和 Metadata Filter。
9. 使用 Rerank 重新排序候选，对比重排前后效果和耗时。
10. 去重、扩展父块、控制 Token，并返回文件名、页码和证据。
11. 使用 MySQL 管理导入状态，实现更新、删除、重试和补偿。
12. 建立测试集，比较切块、Hybrid、Rerank 的 Recall@K、MRR、延迟和成本。

## 完成标志

- 能独立完成文档导入、索引、检索、重排和引用回答。
- 能实现 Parent-Child、Hybrid Search、RRF 和 Rerank。
- 能解释 HNSW 参数和质量/延迟取舍。
- 能维护索引更新、删除、权限和模型升级。
- 能用评估数据证明优化是否有效。

---

# 阶段 4：Tool Calling

## 阶段目标

让模型在受控范围内调用真实业务能力，并保证参数、权限、副作用和审计安全。

## 必学理论

- 模型生成 Tool Call，应用程序负责校验和执行。
- Tool Schema、描述质量和多工具选择。
- 模型参数是不可信输入。
- 身份、租户、权限、幂等和副作用。
- 超时、重试、熔断、二次确认和审计。

## 练习顺序

1. 定义一个本地只读工具，使用 Eino Tool/ToolsNode 执行。
2. 校验必填、类型、范围和未知字段。
3. 将用户和租户身份沿 Context 传递并再次授权。
4. 包装商品、库存、订单和物流等真实业务查询。
5. 注册多个工具，处理不调用、单次调用和连续调用。
6. 处理超时、下游失败、返回过大和有限重试。
7. 实现取消订单/退款模拟，加入确认、幂等和审计。

## 完成标志

- 能将真实业务接口包装为安全工具。
- 能解释参数校验、权限、幂等和重试风险。
- 写操作具备确认、审计和失败处理。

---

# 阶段 5：Agent 工作流

## 阶段目标

使用 Chain、Graph 和 ReAct Agent 完成多步骤任务，同时控制状态、循环、预算和副作用。

## 必学理论

- Chat、RAG、Tool Calling 和 Agent 的关系。
- Chain、Graph、ReAct 与普通状态机的选型。
- 状态、条件分支、循环、终止条件和最大步骤。
- Checkpoint、Interrupt、Resume 和 Human-in-the-loop。
- 固定业务规则优先确定性代码，不交给模型自由决定。

## 练习顺序

1. 使用 Chain 完成固定顺序流程。
2. 使用 Graph 将知识、商品和订单问题路由到不同节点。
3. 实现 ReAct：ChatModel 选择工具，ToolsNode 执行并返回结果。
4. 保存当前步骤、工具结果、错误、Token 和预算状态。
5. 设置最大步骤，检测重复调用和死循环。
6. 在高风险节点 Interrupt，Checkpoint 保存状态，确认后 Resume。
7. 将长任务改为 API 创建任务、Worker 执行、数据库保存状态。
8. 重启后恢复任务，并通过幂等避免重复副作用。

## 完成标志

- 能组合 Retriever 和多个 Tool 构建工作流。
- 能控制最大步骤、失败、成本和人工确认。
- 能持久化、恢复并排查 Agent 执行轨迹。

---

# 阶段 6：Eval 评估

## 阶段目标

用可重复的数据证明 Prompt、模型、检索和 Agent 修改是否真的有效。

## 必学理论

- Golden Dataset 和测试集覆盖范围。
- 检索与回答必须分层评估。
- Recall@K、MRR、命中率、引用正确性和幻觉。
- 工具选择率、参数正确率和 Agent 任务完成率。
- LLM-as-a-Judge 的偏差和人工抽样复核。
- 质量、延迟和成本的共同取舍。

## 练习顺序

1. 建立问题、目标 Chunk、期望事实和期望工具测试集。
2. 自动计算 Retriever 的 Recall@K、MRR 和命中率。
3. 检查回答事实、引用和无依据拒答。
4. 统计 Tool 参数正确率和 Agent 任务完成率。
5. 记录 P50/P95、Token、失败率和成本。
6. 对模型、Prompt、切块、Hybrid 和 Rerank 做回归对比。
7. 将核心评估接入 CI，质量明显下降时阻止发布。

## 完成标志

- 能建立可重复的评估集。
- 能分别评估检索、回答、工具和 Agent。
- 能比较优化前后的质量、性能和成本。

---

# 阶段 7：安全、监控与部署

## 阶段目标

让 AI 系统具备可靠权限、可观测性、故障治理和部署回滚能力。

## 必学理论

- Prompt Injection、间接注入、不可信文档和工具越权。
- Secret、PII、日志脱敏和最小权限。
- Logs、Metrics、Traces 和 SLI/SLO。
- 限流、熔断、Fallback、健康检查、灰度和回滚。

## 练习顺序

1. 使用环境变量/Secret 管理 Key，日志隐藏敏感信息。
2. 构造 Prompt Injection，验证 RAG 和 Tool 权限边界。
3. 为接口、模型和 Embedding 增加限流与并发控制。
4. 使用 trace ID 串联模型、Retriever、Rerank、Tool 和 Agent。
5. 记录耗时、Token、成本、错误率和队列长度并配置告警。
6. 使用 Docker Compose 启动应用、MySQL、Redis 和向量库。
7. 配置健康检查、灰度发布和回滚。
8. 模拟主模型或向量库不可用并验证降级。

## 完成标志

- 权限在检索和工具执行阶段生效。
- 能追踪一条完整 AI 请求链路。
- 能处理限流、故障、部署和回滚。

---

# 进阶专题：按岗位或项目需要学习

以下内容知道用途即可，不要求在主线阶段逐个写 Demo：

- 本地模型、Hugging Face、Fine-tuning、LoRA 和多模态。
- Semantic Chunking、Multi-Vector、Late Chunking、HyDE。
- RAPTOR、GraphRAG、知识图谱和 ColBERT。
- IVF_PQ 等向量压缩索引细节。
- MCP、复杂 Tool Registry 和完整 Saga 平台。
- Plan-and-Execute、Supervisor、多 Agent 和长期记忆框架。
- Kubernetes、Service Mesh、多云模型网关和跨区域容灾。

当岗位 JD 或实际项目明确需要某项时，再单独建立专题学习文档。

---

# 推荐学习时间

如果每天学习 1～2 小时：

| 阶段 | 建议时间 |
| --- | --- |
| Go 后端基础 | 1～2 个月 |
| LLM API 与 Prompt | 1～2 个月 |
| RAG | 3～4 个月 |
| Tool Calling | 1～2 个月 |
| Agent | 1～2 个月 |
| Eval | 1 个月 |
| 安全、监控与部署 | 1～2 个月 |

整体约 10～15 个月。学习进度以阶段完成标志为准，不以时间为准。

---

# 官方资料

- [CloudWeGo Eino 官方文档](https://www.cloudwego.io/zh/docs/eino/)
- [Eino 核心组件](https://www.cloudwego.io/zh/docs/eino/core_modules/)
- [Eino Examples](https://github.com/cloudwego/eino-examples)
- [Eino Checkpoint、Interrupt 与 Resume](https://www.cloudwego.io/zh/docs/eino/core_modules/chain_and_graph_orchestration/checkpoint_interrupt)
- [Eino ReAct Agent](https://www.cloudwego.io/zh/docs/eino/core_modules/flow_integration_components/react_agent_manual)

具体练习开始前，再查对应组件的当前官方文档和最小示例。
