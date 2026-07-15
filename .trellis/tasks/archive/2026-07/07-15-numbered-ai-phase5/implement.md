# 阶段 5：Agent 工作流编号练习执行计划

## 实现清单

- [x] 生成 `01_chain_workflow`：固定顺序 Chain。
- [x] 生成 `02_graph_routing`：知识、商品与订单 Graph 路由。
- [x] 生成 `03_react_agent`：ChatModel、ToolsNode 与 ReAct。
- [x] 生成 `04_agent_state_and_budget`：步骤、结果、错误、Token 与预算状态。
- [x] 生成 `05_loop_detection`：最大步骤、重复调用与死循环检测。
- [x] 生成 `06_interrupt_checkpoint_resume`：Interrupt、Checkpoint 与 Resume。
- [x] 生成 `07_async_agent_task`：API 创建任务、Worker 执行与状态持久化。
- [x] 生成 `08_task_recovery_idempotency`：重启恢复与幂等副作用。

## 一致性检查

- [x] 验证目录编号从 `01` 连续到 `08`。
- [x] 验证每个目录包含 README，章节和 TODO 编号连续。
- [x] 验证默认配置不会连接外部服务或产生付费请求。
- [x] 搜索本阶段旧路径、重复目录、真实密钥和未处理错误。

## 验证

```bash
gofmt -w <阶段 5 新增或迁移的 Go 文件>
go test -timeout=60s ./examples/ai/phase5/...
go vet ./examples/ai/phase5/...
git diff --check
```

阶段 5 完成后执行 `trellis-check` 和 `go-review`；涉及 SQL/Schema/查询时追加 `sql-code-review`。
