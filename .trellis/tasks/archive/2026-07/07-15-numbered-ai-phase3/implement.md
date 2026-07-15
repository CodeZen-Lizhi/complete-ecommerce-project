# 阶段 3：RAG 编号练习执行计划

## 实现清单

- [x] 生成 `01_embedding_similarity`：Embedding、余弦相似度与 Top-K。
- [x] 生成 `02_document_loading_metadata`：Markdown/PDF 加载与 Metadata。
- [x] 生成 `03_chunking_strategies`：递归、结构化、Parent-Child 与表格切块。
- [x] 生成 `04_vector_index_and_retrieval`：Eino Indexer 与 Retriever。
- [x] 生成 `05_hnsw_tuning`：HNSW 参数、命中率与 P95。
- [x] 生成 `06_parent_child_retrieval`：子块召回与父块返回。
- [x] 生成 `07_hybrid_search_rrf`：Dense、BM25 与 RRF。
- [x] 生成 `08_query_rewrite_and_filter`：多轮改写、Multi-Query 与 Metadata Filter。
- [x] 生成 `09_rerank`：候选重排与耗时对比。
- [x] 生成 `10_context_budget_and_citations`：去重、父块扩展、Token 与引用。
- [x] 生成 `11_index_lifecycle`：导入状态、更新、删除、重试与补偿。
- [x] 生成 `12_rag_evaluation`：Recall@K、MRR、延迟与成本。

## 一致性检查

- [x] 验证目录编号从 `01` 连续到 `12`。
- [x] 验证每个目录包含 README，章节和 TODO 编号连续。
- [x] 验证默认配置不会连接外部服务或产生付费请求。
- [x] 搜索本阶段旧路径、重复目录、真实密钥和未处理错误。

## 验证

```bash
gofmt -w <阶段 3 新增或迁移的 Go 文件>
go test -timeout=60s ./examples/ai/phase3/...
go vet ./examples/ai/phase3/...
git diff --check
```

阶段 3 完成后执行 `trellis-check` 和 `go-review`；涉及 SQL/Schema/查询时追加 `sql-code-review`。
