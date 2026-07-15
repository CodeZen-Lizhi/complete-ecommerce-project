# AI 学习练习索引

本目录按学习路线分为 7 个阶段、共 60 个练习。每个阶段从 `01` 重新编号，建议严格按顺序完成。

- 阶段 1 是现有电商项目的导航型练习：README 指向实际修改位置。
- 阶段 2–7 是隔离型学习骨架：核心实现保留为连续中文 TODO。
- 默认占位配置不会连接外部模型、数据库、Redis 或向量库。

## 阶段 1：Go 后端基础

- [01_gin_validated_api](./phase1/01_gin_validated_api/README.md)：Gin 参数校验与统一响应
- [02_gorm_crud_pagination_transaction](./phase1/02_gorm_crud_pagination_transaction/README.md)：GORM CRUD、分页、事务与 Explain
- [03_redis_cache_aside](./phase1/03_redis_cache_aside/README.md)：Redis Cache Aside
- [04_jwt_auth_and_authorization](./phase1/04_jwt_auth_and_authorization/README.md)：JWT 鉴权与授权
- [05_http_client_resilience](./phase1/05_http_client_resilience/README.md)：HTTP Client 韧性
- [06_worker_pool](./phase1/06_worker_pool/README.md)：Worker Pool
- [07_backend_tests](./phase1/07_backend_tests/README.md)：后端测试

## 阶段 2：LLM API 与 Prompt

- [01_basic_chat](./phase2/01_basic_chat/README.md)：基础模型调用
- [02_eino_generate](./phase2/02_eino_generate/README.md)：Eino Generate
- [03_prompt_roles](./phase2/03_prompt_roles/README.md)：System 与 User Prompt
- [04_in_memory_multi_turn](./phase2/04_in_memory_multi_turn/README.md)：程序内多轮对话
- [05_redis_session_history](./phase2/05_redis_session_history/README.md)：Redis 会话历史
- [06_session_history_limits](./phase2/06_session_history_limits/README.md)：历史截断、TTL 与隔离
- [07_chat_template](./phase2/07_chat_template/README.md)：ChatTemplate
- [08_streaming_chat](./phase2/08_streaming_chat/README.md)：Stream
- [09_structured_json_output](./phase2/09_structured_json_output/README.md)：结构化 JSON
- [10_call_governance](./phase2/10_call_governance/README.md)：调用治理
- [11_model_provider_adapter](./phase2/11_model_provider_adapter/README.md)：模型适配层

## 阶段 3：RAG

- [01_embedding_similarity](./phase3/01_embedding_similarity/README.md)：Embedding 与 Top-K
- [02_document_loading_metadata](./phase3/02_document_loading_metadata/README.md)：文档加载与 Metadata
- [03_chunking_strategies](./phase3/03_chunking_strategies/README.md)：切块策略
- [04_vector_index_and_retrieval](./phase3/04_vector_index_and_retrieval/README.md)：索引与检索
- [05_hnsw_tuning](./phase3/05_hnsw_tuning/README.md)：HNSW 调优
- [06_parent_child_retrieval](./phase3/06_parent_child_retrieval/README.md)：Parent-Child
- [07_hybrid_search_rrf](./phase3/07_hybrid_search_rrf/README.md)：Hybrid 与 RRF
- [08_query_rewrite_and_filter](./phase3/08_query_rewrite_and_filter/README.md)：改写与过滤
- [09_rerank](./phase3/09_rerank/README.md)：Rerank
- [10_context_budget_and_citations](./phase3/10_context_budget_and_citations/README.md)：预算与引用
- [11_index_lifecycle](./phase3/11_index_lifecycle/README.md)：索引生命周期
- [12_rag_evaluation](./phase3/12_rag_evaluation/README.md)：RAG 评估

## 阶段 4：Tool Calling

- [01_local_readonly_tool](./phase4/01_local_readonly_tool/README.md)：本地只读工具
- [02_tool_argument_validation](./phase4/02_tool_argument_validation/README.md)：参数校验
- [03_context_identity_authorization](./phase4/03_context_identity_authorization/README.md)：身份与授权
- [04_ecommerce_query_tools](./phase4/04_ecommerce_query_tools/README.md)：电商查询工具
- [05_multi_tool_orchestration](./phase4/05_multi_tool_orchestration/README.md)：多工具编排
- [06_tool_resilience](./phase4/06_tool_resilience/README.md)：失败治理
- [07_safe_write_tool](./phase4/07_safe_write_tool/README.md)：安全写工具

## 阶段 5：Agent 工作流

- [01_chain_workflow](./phase5/01_chain_workflow/README.md)：Chain
- [02_graph_routing](./phase5/02_graph_routing/README.md)：Graph 路由
- [03_react_agent](./phase5/03_react_agent/README.md)：ReAct
- [04_agent_state_and_budget](./phase5/04_agent_state_and_budget/README.md)：状态与预算
- [05_loop_detection](./phase5/05_loop_detection/README.md)：循环检测
- [06_interrupt_checkpoint_resume](./phase5/06_interrupt_checkpoint_resume/README.md)：Interrupt/Checkpoint/Resume
- [07_async_agent_task](./phase5/07_async_agent_task/README.md)：异步任务
- [08_task_recovery_idempotency](./phase5/08_task_recovery_idempotency/README.md)：恢复与幂等

## 阶段 6：Eval

- [01_golden_dataset](./phase6/01_golden_dataset/README.md)：Golden Dataset
- [02_retrieval_metrics](./phase6/02_retrieval_metrics/README.md)：检索指标
- [03_answer_quality_checks](./phase6/03_answer_quality_checks/README.md)：回答质量
- [04_tool_and_agent_metrics](./phase6/04_tool_and_agent_metrics/README.md)：Tool 与 Agent 指标
- [05_latency_token_cost_metrics](./phase6/05_latency_token_cost_metrics/README.md)：延迟 Token 成本
- [06_regression_comparison](./phase6/06_regression_comparison/README.md)：回归对比
- [07_eval_ci_gate](./phase6/07_eval_ci_gate/README.md)：CI 门禁

## 阶段 7：安全、监控与部署

- [01_secret_and_log_redaction](./phase7/01_secret_and_log_redaction/README.md)：Secret 与脱敏
- [02_prompt_injection_defense](./phase7/02_prompt_injection_defense/README.md)：Prompt Injection
- [03_rate_limit_and_concurrency](./phase7/03_rate_limit_and_concurrency/README.md)：限流与并发
- [04_end_to_end_tracing](./phase7/04_end_to_end_tracing/README.md)：端到端 Trace
- [05_metrics_and_alerts](./phase7/05_metrics_and_alerts/README.md)：Metrics 与告警
- [06_docker_compose_stack](./phase7/06_docker_compose_stack/README.md)：Docker Compose
- [07_health_gray_release_rollback](./phase7/07_health_gray_release_rollback/README.md)：健康、灰度与回滚
- [08_dependency_failure_fallback](./phase7/08_dependency_failure_fallback/README.md)：故障与降级

## 通用验证

```bash
go test -timeout=60s ./examples/ai/...
go vet ./examples/ai/...
```

阶段 1 只有导航 README；完成其中练习后，应按各目录说明验证实际电商代码。
