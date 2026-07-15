# 全阶段编号练习执行计划

## 1. 规划与任务树

- [x] 创建阶段 1 至阶段 7 子任务，并在各子任务 PRD 中列出目录映射和验收标准。
- [x] 为每个阶段补充必要设计、实现顺序和官方资料研究入口。
- [x] 先完成父任务目录契约和全量索引设计，再逐阶段实施。

## 2. 阶段实施顺序

- [x] 阶段 1：生成 7 个编号导航目录，链接现有 Handler、Service、Repository、Redis、JWT、HTTP Client、并发和测试练习位置。
- [x] 阶段 2：将已有 6 个目录迁移为编号目录，更新全部引用；生成第 7 至第 11 个可编译骨架。
- [x] 阶段 3：生成 12 个 RAG 骨架，按 Embedding、文档处理、索引、检索、重排、生命周期和评估递进。
- [x] 阶段 4：生成 7 个 Tool Calling 骨架，覆盖参数校验、身份授权、真实业务查询、多工具、失败治理和写操作安全。
- [x] 阶段 5：生成 8 个 Agent 骨架，覆盖 Chain、Graph、ReAct、状态、循环控制、人工确认、异步任务和恢复。
- [x] 阶段 6：生成 7 个 Eval 骨架，覆盖数据集、检索、回答、工具、Agent、性能成本和 CI 回归。
- [x] 阶段 7：生成 8 个安全、监控与部署骨架，覆盖 Secret、注入、限流、追踪、指标、Compose、灰度回滚和故障降级。

## 3. 父任务集成

- [x] 在学习路线文档中增加每项练习到编号目录的映射。
- [x] 搜索并替换阶段 2 旧目录引用，确认不存在重复目录或失效命令。
- [x] 检查所有 README 的章节、编号、环境变量和“暂不实现”边界一致。
- [x] 检查所有 Go TODO 在单文件内连续编号，未完成路径不会连接外部服务。

## 4. 验证命令

```bash
gofmt -w <本任务新增或迁移的 Go 文件>
go test -timeout=60s ./...
go vet ./...
git diff --check
rg -n 'examples/ai/phase2/(basic_chat|eino_generate|prompt_roles|in_memory_multi_turn|redis_session_history|session_history_limits)' . \
  --glob '!.trellis/tasks/archive/**'
```

按阶段追加：

- README/目录完整性脚本：验证 7 个阶段、60 个编号目录、连续序号和必需文件。
- 外部服务练习：只验证默认占位配置的安全退出，不在自动测试中连接真实服务。
- 涉及 SQL 的练习：执行 SQL 专项静态 review；没有真实 Schema 时不伪造迁移验证。

## 5. Review 与回滚点

- [x] 每个阶段完成后执行 `go-review`，当前范围内明确问题直接修复并重新验证。
- [x] 涉及 SQL/Schema/查询的阶段追加 `sql-code-review`。
- [x] 父任务最后执行全仓 Go review、目录一致性检查和旧路径扫描。
- [x] 每个阶段保持独立变更边界；目录迁移失败时恢复该阶段路径，不使用破坏性 Git 命令。
