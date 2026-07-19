# 补齐 AI 学习路线缺失练习与行为验收：实施计划

## 1. 阶段 2 Few-shot

- [x] 在 Prompt 角色练习后新增 Zero-shot/Few-shot 对比目录。
- [x] 提供真实 AgenticModel 配置、消息构造函数签名、分步骤 TODO 和占位 Key 校验入口。
- [x] 顺延阶段 2 后续目录编号并修正相互引用。
- [x] 更新路线、总索引和 `catalog_test.go` 的目录数量/专属符号契约。
- [x] 验证阶段 2 所有目录编译。

## 2. 阶段 4 真实 Tool Calling

- [x] 加强本地只读工具练习，明确使用 Eino Tool 和 `compose.NewToolNode`。
- [x] 加强多工具编排练习，要求真实 ChatModel Tool Call、ToolCallID 对应和工具结果回传模型。
- [x] 删除“等价 toolExecutor 即可完成”的绕行表述。
- [x] 补充真实运行命令、输出观察项和错误场景。

## 3. 阶段 5 真实 Eino 编排

- [x] Chain 骨架改为真实 `compose.NewChain`、节点注册、Compile 和 Invoke 承接。
- [x] Graph 骨架改为真实 `compose.NewGraph`、分支边、Compile 和 Invoke 承接。
- [x] ReAct 骨架改为真实 `react.NewAgent` 与 `compose.ToolsNodeConfig`。
- [x] Interrupt/Checkpoint/Resume 骨架改为仓库 Eino 版本对应的真实 API。
- [x] README 明确每个真实组件的关键类型、调用顺序和完成标准。

## 4. 阶段 3 真实 RAG 基础设施

- [x] PDF Loader 增加真实 Adapter 构造点和页级 fixture 运行路径。
- [x] 向量索引、HNSW、Hybrid、Rerank 练习明确默认真实实现和配置。
- [x] Qdrant 路径与阶段 7 Compose 配置保持一致。
- [x] MySQL 索引生命周期增加真实状态存储构造入口、示例表契约和恢复运行步骤。

## 5. 阶段 6～7 真实评估与可观测性

- [x] LLM-as-a-Judge 增加真实模型、结构化 Schema、人工标签对比 TODO。
- [x] Tracing 增加 OpenTelemetry TracerProvider/Exporter/Span 骨架。
- [x] 灰度与回滚增加真实版本路由和演练入口。
- [x] Docker Compose 按需要补充可观测组件，但不扩展到生产 Kubernetes（本次采用 stdout Exporter，无需新增服务）。

## 6. 一致性扫描

- [x] 全量搜索“暂不实现”“后续练习”“等价接口”等表述，确认都有真实承接或明确范围。
- [x] 对照路线知识点、练习 TODO、README 完成标准和总索引。
- [x] 确认未新增练习目录内 `_test.go`，未写入真实凭证。

## 验证命令

```bash
gofmt -w <本次修改的 Go 文件>
go test -timeout=60s ./...
go vet ./...
```

按批次追加 README 中声明的真实运行命令；缺少用户本地服务或模型配置时，记录未执行原因，不使用 Fake 结果替代。

## Review 门禁

- 主 agent 使用 `go-review` 检查所有 Go 改动。
- 涉及 MySQL 表契约、查询或索引时追加 `sql-code-review`。
- 检查是否仍存在可绕过真实 Eino/基础设施调用的完成路径。
- 检查新函数中文注释、Context 传播、错误链和敏感配置处理。

## 回滚点

- 批次 1：阶段 2 新目录和编号调整。
- 批次 2：阶段 4 Tool Calling。
- 批次 3：阶段 5 Eino 编排。
- 批次 4：阶段 3 RAG 基础设施。
- 批次 5：阶段 6～7 评估与可观测性。
