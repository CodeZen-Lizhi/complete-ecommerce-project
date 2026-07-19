# AI 学习练习规范

## Scenario: 真实组件 TODO 骨架

### 1. Scope / Trigger

- 适用于 `examples/ai/` 中声称教授 Eino、模型服务、向量库、数据库或可观测组件的练习。
- 目标是避免路线写了真实知识点，但练习只留下可替代的本地接口，导致学习者绕过目标 API。

### 2. Signatures

- 骨架应出现目标 SDK 的真实类型或构造入口，例如：
  - `compose.NewToolNode` / `*compose.ToolsNode`
  - `compose.NewChain` / `compose.NewGraph` / `compose.Runnable`
  - `react.NewAgent` / `*react.Agent`
  - `compose.CheckPointStore` / `compose.ResumeWithData`
- `runExercise(ctx context.Context) error` 继续作为练习执行入口，未完成时返回 `errExerciseIncomplete`。
- 编号 TODO 必须贴在学习者实际填写的类型、构造函数、校验函数或调用位置；`runExercise` 只负责按顺序串联步骤，不得集中罗列整组 TODO。

### 3. Contracts

- 模型配置：需要真实模型调用的练习在 `exercise.go` 或 `main.go` 顶部提供 `BaseURL`、占位 `API Key` 和模型名常量，学习者直接修改这些入口；真实 Key 不提交。
- 基础设施配置：Qdrant、BM25、MySQL、Reranker、发布控制等练习也在顶部提供服务专属占位常量，README 必须列出应修改的常量名。
- Secret 管理与 Docker Compose 练习是例外：它们的学习目标就是环境变量和 Secret 生命周期，不得把真实密钥改成源码常量。
- README 必须包含真实运行命令、是否产生模型费用、预期观察项和失败场景。
- 用户明确要求时，练习保持 TODO 骨架，不直接交付核心答案。

### 4. Validation & Error Matrix

| 条件 | 行为 |
| --- | --- |
| Context 为 nil | 返回明确错误 |
| API Key 为空或仍为占位值 | 真实调用前失败 |
| 外部服务地址为空 | 构造 Client 前失败 |
| SDK 构造或真实调用失败 | 使用 `%w` 保留底层错误链 |
| 只实现 Fake/Mock/等价接口 | 不满足 README 完成标准 |

### 5. Good/Base/Bad Cases

- Good：骨架保留小接口用于理解边界，同时提供真实 SDK 构造函数 TODO，并要求运行真实组件闭环。
- Base：目录、README、Go 骨架可编译，未完成时明确返回 `errExerciseIncomplete`。
- Bad：README 写“Eino 或等价实现”，最终只调用本地 Map、Fake 模型或手写状态机。
- Bad：多个编号 TODO 连续堆在入口函数里，但对应类型和函数没有可填写位置。

### 6. Tests Required

- 仓库级 `catalog_test.go` 验证目录、TODO 连续性、TODO 填写位置、README 一致性、顶部占位配置和真实构造入口符号存在。
- `go test -timeout=60s ./...` 和 `go vet ./...` 保证骨架基线。
- 用户要求不新增 `_test.go` 时，真实完成证明来自 README 指定的 `go run`、Compose 或服务调用命令，不用 Fake 结果替代。

### 7. Wrong vs Correct

#### Wrong

```go
// TODO：创建 Eino ReAct Agent 或等价本地循环。
func runReAct(model localModel, tools localTools) error
```

#### Correct

```go
// TODO：通过 react.NewAgent 和 compose.ToolsNodeConfig 创建真实 Eino Agent。
func newEinoReActAgent(
    ctx context.Context,
    chatModel model.ToolCallingChatModel,
    tools []tool.BaseTool,
    maxSteps int,
) (*react.Agent, error)
```
