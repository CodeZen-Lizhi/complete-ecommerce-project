# 补齐 AI 学习路线缺失练习与行为验收：技术设计

## 设计目标

让学习路线中的关键知识点必须通过真实 SDK、真实服务或真实模型调用完成，同时保持练习代码是待填写骨架，而不是现成答案。

## 总体结构

### 新增练习

- 在阶段 2 新增独立 Few-shot 练习，放在现有 Prompt 角色练习之后。
- 后续阶段 2 练习顺延编号，并同步总索引、路线和目录契约计数。
- 新练习沿用阶段 2 已有结构：`README.md`、`main.go`，配置从已有占位常量和环境约定读取，不写真实 Key。

### 加强既有练习

- 阶段 3～7 保留现有目录和编号。
- 删除 README/TODO 中允许使用“等价自定义实现”作为最终完成方式的表述。
- 自定义接口可以保留为业务边界，但必须新增明确的真实 Adapter 构造点、真实组件类型或编译期类型断言。
- `runExercise` 继续以 `errExerciseIncomplete` 作为未完成提示；关键 TODO 要求学习者完成真实连接和调用。

## 真实承接边界

### Eino

- 基于仓库锁定的 `github.com/cloudwego/eino v0.9.12` 设计骨架。
- Tool 使用 `components/tool` 与 `compose.NewToolNode`。
- Chain 使用 `compose.NewChain[I, O]` 并 `Compile`。
- Graph 使用 `compose.NewGraph[I, O]`、节点/边注册和 `Compile`。
- ReAct 使用 `flow/agent/react.NewAgent`，真实绑定 ChatModel 和 `compose.ToolsNodeConfig`。
- Interrupt/Resume 使用当前版本真实 interrupt state、checkpoint 或 runner resume API；以本地模块源码和官方文档交叉确认签名。

### RAG 基础设施

- PDF Loader：骨架必须包含真实页级 Loader Adapter 构造点和本地 fixture 路径，不只保留通用接口。
- 向量库：默认实践路径统一到现有 Compose 中的 Qdrant，配置从 `VECTOR_BASE_URL` 等环境变量读取。
- MySQL：索引生命周期练习必须提供真实状态存储 Adapter 的构造入口和 DSN 配置，不臆造生产 Schema；练习所需表契约放在示例范围内说明。
- HNSW：通过真实向量库索引参数执行多组实验，不以纯内存排序模拟完成。

### 评估与可观测性

- LLM-as-a-Judge 使用真实 ChatModel，并要求结构化 JSON Schema 返回评分、理由和证据；人工标签对比仍是必做步骤。
- Tracing 使用 OpenTelemetry API 创建父子 Span，并提供本地 Collector/可观察后端的运行说明。
- 灰度练习保留决策逻辑骨架，并增加真实版本标识、路由比例和回滚场景入口。

## 配置与安全

- 沿用仓库学习示例的占位 Key 约定；不得提交真实 API Key。
- 外部基础设施地址通过环境变量或阶段 7 Compose 获取。
- 运行真实模型前必须校验占位 Key，失败时明确退出。
- README 明确哪些命令会产生模型调用或费用。

## 验收方式

- 不新增练习目录内 `_test.go`。
- README 为每个练习提供：启动依赖、运行命令、待观察输出、至少一个失败场景、完成标准。
- 仓库级 `go test` 只保证骨架、目录契约和已有代码可编译，不作为真实练习完成证明。
- 真实练习的完成证明来自学习者运行 `go run` 或 README 指定命令后获得的真实服务/模型结果。

## 兼容与回滚

- 保留现有用户改动，尤其是阶段 2 练习 9 和 `go.mod`。
- 分批修改，每一批独立通过格式化、测试和 vet；若某批 SDK 签名不兼容，可单独回滚该批而不影响其他练习。
- 不修改生产商城代码和生产数据库结构。
