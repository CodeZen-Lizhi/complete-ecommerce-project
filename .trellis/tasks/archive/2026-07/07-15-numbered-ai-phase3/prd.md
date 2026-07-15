# 阶段 3：RAG 编号练习

## Goal

生成 12 个递进式 RAG 练习骨架，覆盖 Embedding、文档处理、切块、索引、检索、HNSW、Parent-Child、Hybrid、改写、重排、引用、索引生命周期和评估。

## Parent Task

- 父任务：`.trellis/tasks/07-15-numbered-ai-exercises/`
- 本子任务只负责阶段 3 的目录、骨架、文档和阶段内验证；跨阶段索引与最终全仓检查由父任务负责。

## Exercise Map

| 顺序 | 目录 | 主题 |
| --- | --- | --- |
| 1 | `01_embedding_similarity` | Embedding、余弦相似度与 Top-K |
| 2 | `02_document_loading_metadata` | Markdown/PDF 加载与 Metadata |
| 3 | `03_chunking_strategies` | 递归、结构化、Parent-Child 与表格切块 |
| 4 | `04_vector_index_and_retrieval` | Eino Indexer 与 Retriever |
| 5 | `05_hnsw_tuning` | HNSW 参数、命中率与 P95 |
| 6 | `06_parent_child_retrieval` | 子块召回与父块返回 |
| 7 | `07_hybrid_search_rrf` | Dense、BM25 与 RRF |
| 8 | `08_query_rewrite_and_filter` | 多轮改写、Multi-Query 与 Metadata Filter |
| 9 | `09_rerank` | 候选重排与耗时对比 |
| 10 | `10_context_budget_and_citations` | 去重、父块扩展、Token 与引用 |
| 11 | `11_index_lifecycle` | 导入状态、更新、删除、重试与补偿 |
| 12 | `12_rag_evaluation` | Recall@K、MRR、延迟与成本 |

## Requirements

- 阶段内目录使用两位数字连续编号，并与路线图顺序一致。
- 独立 Go 骨架与小型本地数据集；外部向量库、模型和 PDF 解析依赖通过配置或接口隔离。
- 每个练习至少包含中文 README，结构包含目标、前置知识、连续 TODO、运行/验证、完成标准和暂不实现。
- 所有新增 Go 函数和方法写简洁中文职责注释；核心 TODO 按真实执行顺序连续编号。
- 默认占位配置不得连接外部服务、泄露真实密钥或修改外部状态。
- 第一题吸收现有 Embedding 规划中的 CLI、环境变量和 Eino Embedder 契约；旧的 `.trellis/tasks/07-13-ai-labs-embedding/` 保持为本地未跟踪任务，不纳入本次提交，避免生成第二套练习目录。

## Acceptance Criteria

- [x] 阶段 3 的 12 个编号目录全部存在，序号连续且没有重复目录。
- [x] 每个 README 的主题、TODO 和完成标准与路线图对应项一致。
- [x] 需要编码的骨架可以编译；未完成路径明确失败且默认不访问外部服务。
- [x] 阶段内 Go 文件通过 `gofmt`、目标包测试与 `go vet`。
- [x] 使用 `go-review` 完成阶段 review；涉及 SQL/Schema/查询时追加 `sql-code-review`。

## Out of Scope

- 提供完整参考答案或生产级实现。
- 修改其他阶段练习或现有电商业务语义。
- 使用真实云资源、真实凭证或执行不可逆外部操作。
