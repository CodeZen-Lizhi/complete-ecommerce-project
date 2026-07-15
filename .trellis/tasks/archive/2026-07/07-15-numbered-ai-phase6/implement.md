# 阶段 6：Eval 编号练习执行计划

## 实现清单

- [x] 生成 `01_golden_dataset`：问题、目标 Chunk、事实与工具测试集。
- [x] 生成 `02_retrieval_metrics`：Recall@K、MRR 与命中率。
- [x] 生成 `03_answer_quality_checks`：事实、引用与无依据拒答。
- [x] 生成 `04_tool_and_agent_metrics`：Tool 参数正确率与 Agent 完成率。
- [x] 生成 `05_latency_token_cost_metrics`：P50/P95、Token、失败率与成本。
- [x] 生成 `06_regression_comparison`：模型、Prompt、切块、Hybrid 与 Rerank 对比。
- [x] 生成 `07_eval_ci_gate`：核心评估接入 CI 与质量门禁。

## 一致性检查

- [x] 验证目录编号从 `01` 连续到 `07`。
- [x] 验证每个目录包含 README，章节和 TODO 编号连续。
- [x] 验证默认配置不会连接外部服务或产生付费请求。
- [x] 搜索本阶段旧路径、重复目录、真实密钥和未处理错误。

## 验证

```bash
gofmt -w <阶段 6 新增或迁移的 Go 文件>
go test -timeout=60s ./examples/ai/phase6/...
go vet ./examples/ai/phase6/...
git diff --check
```

阶段 6 完成后执行 `trellis-check` 和 `go-review`；涉及 SQL/Schema/查询时追加 `sql-code-review`。
