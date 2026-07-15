# 阶段 7：安全监控与部署编号练习执行计划

## 实现清单

- [x] 生成 `01_secret_and_log_redaction`：Secret 管理与日志脱敏。
- [x] 生成 `02_prompt_injection_defense`：Prompt Injection 与权限边界。
- [x] 生成 `03_rate_limit_and_concurrency`：接口、模型与 Embedding 限流。
- [x] 生成 `04_end_to_end_tracing`：模型、Retriever、Rerank、Tool 与 Agent Trace。
- [x] 生成 `05_metrics_and_alerts`：耗时、Token、成本、错误率与队列告警。
- [x] 生成 `06_docker_compose_stack`：应用、MySQL、Redis 与向量库 Compose。
- [x] 生成 `07_health_gray_release_rollback`：健康检查、灰度与回滚。
- [x] 生成 `08_dependency_failure_fallback`：模型或向量库故障与降级。

## 一致性检查

- [x] 验证目录编号从 `01` 连续到 `08`。
- [x] 验证每个目录包含 README，章节和 TODO 编号连续。
- [x] 验证默认配置不会连接外部服务或产生付费请求。
- [x] 搜索本阶段旧路径、重复目录、真实密钥和未处理错误。

## 验证

```bash
gofmt -w <阶段 7 新增或迁移的 Go 文件>
go test -timeout=60s ./examples/ai/phase7/...
go vet ./examples/ai/phase7/...
git diff --check
```

阶段 7 完成后执行 `trellis-check` 和 `go-review`；涉及 SQL/Schema/查询时追加 `sql-code-review`。
