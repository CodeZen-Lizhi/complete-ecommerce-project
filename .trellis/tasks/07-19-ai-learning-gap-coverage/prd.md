# 补齐 AI 学习路线缺失练习与行为验收

## Goal

补齐 AI 学习路线中“文档提到但练习没有真实覆盖”以及“只有接口骨架、没有真实调用”的缺口，使学习者必须在可执行练习中实际使用目标 API 或基础设施，并能通过运行结果判断是否完成。

## Background

- 阶段 2 的路线包含 Few-shot，但当前 Prompt 练习只对比 System Prompt。
- 阶段 3～7 当前主要是抽象接口和 TODO；用户明确不要求为这些练习新增 `_test.go`，而要求补足真实运行路径。
- Tool Calling、Chain、Graph、ReAct、Checkpoint 等练习存在“Eino 或等价自定义接口”的绕行路径，不能证明学习者使用了真实 Eino API。
- PDF、向量库、MySQL、Tracing 和 LLM-as-a-Judge 多数只有抽象接口或 TODO，没有完整 Adapter、调用路径和可执行验收。
- 阶段 2 练习 9 已补充 Prompt JSON 与原生 JSON Schema 双模式，本任务不得回退或重复实现该工作。

## Requirements

### R1 可执行验收基线

- 不为本任务新增练习目录内的 `_test.go`。
- 每个被补充的知识点必须落实到 `main.go`/`exercise.go` 的真实调用 TODO、可执行入口、运行命令和预期观察结果。
- 完成标准至少覆盖成功路径，以及非法输入、依赖错误、Context 取消等适用失败路径；不得只要求定义接口或类型。

### R2 Few-shot 练习

- 阶段 2 新增可独立完成的 Zero-shot 与 Few-shot 对比练习目录，不把该知识点塞入已有 System Prompt 练习。
- 固定模型、问题和其他参数，只改变示例消息，避免混淆变量。
- 覆盖消息顺序、示例边界、格式约束和示例泄漏风险。
- 新目录只提供待学习者填写的 TODO 骨架；类型、配置入口、调用顺序和观察目标完整，核心模型调用与比较逻辑不直接给出答案。

### R3 真实 Eino API 承接

- Tool/ToolsNode、Chain、Graph、ReAct、Interrupt/Checkpoint/Resume 必须各有真实 Eino API 的实现或明确 Adapter。
- 可保留小接口用于表达边界，但“只实现等价接口”不得满足最终完成标准。
- README 必须指出官方组件、真实调用入口、配置来源和离线测试边界。
- 已有但允许绕行的练习直接加强 TODO、入口和完成标准，不为同一知识点重复创建平行目录。

### R4 真实 RAG 与运行基础设施承接

- PDF Loader、向量索引/检索、HNSW、MySQL 索引生命周期至少提供一条明确、可运行的默认技术路径。
- 外部依赖必须从环境变量或本地 Compose 获取配置，不提交真实凭证。
- 用户明确选择真实集成路径：完成标准必须实际连接目标组件或调用目标模型，不以 Fake、Mock 或纯内存实现替代真实练习。
- 真实调用必须由明确命令触发，并在执行前校验配置，避免因占位配置产生不可诊断结果。

### R5 评估与可观测性承接

- LLM-as-a-Judge 应包含结构化评分契约、真实模型调用和人工标签一致性比较，不得作为唯一质量门禁。
- Tracing 应承接真实 OpenTelemetry API；灰度练习应通过可执行场景覆盖版本路由、最小样本、扩量和回滚决策。

### R6 兼容性与范围控制

- 保留既有练习的相对学习顺序；Few-shot 等完全缺失的知识点新增目录并同步后续编号、路线文档、总索引和目录契约。
- 不修改生产电商业务模块，不把学习示例测试表述成生产覆盖。
- 保留用户当前未提交的练习 9、路线文档和 `go.mod` 改动，不覆盖或回滚。
- 所有新增和加强的练习均采用“待填写 TODO 骨架”形式，不交付核心知识点的完整参考答案。

## Acceptance Criteria

- [ ] 本任务不新增练习目录内 `_test.go`；每个新增或调整的练习都有真实运行入口、命令、预期输出和失败观察项。
- [ ] Few-shot 有独立练习、README、实现骨架和真实模型对比入口。
- [ ] Eino ToolsNode、Chain、Graph、ReAct、Checkpoint/Resume 均存在不能被自定义空接口替代的真实 API 验收点。
- [ ] RAG 默认实践路径说明所用 Loader、向量库和状态存储，并提供不泄露凭证的配置方式。
- [ ] 集成验收实际连接 MySQL、Qdrant、OpenTelemetry 和所选模型服务；Fake/Mock 结果不能作为练习完成证明。
- [ ] LLM-as-a-Judge 和 OpenTelemetry 分别存在真实调用承接及可执行验证步骤。
- [ ] `go test -timeout=60s ./...` 与 `go vet ./...` 继续通过，用于保证目录和代码基线；练习完成证明来自 README 指定的真实运行命令。
- [ ] 所有新增 Go 函数和方法具有符合项目规范的简洁中文注释。
- [ ] README、路线、练习索引、TODO 和测试完成标准彼此一致，不再出现“后续补充但后续没有承接”。

## Out of Scope

- 生产 Kubernetes、云托管向量数据库、生产发布平台和在线 A/B 平台。
- 端到端训练 Embedding 或 Reranker 模型。
- 修改生产商城 Handler、Service、Repository 或数据库 Schema。
