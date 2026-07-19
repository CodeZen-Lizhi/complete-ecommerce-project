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

### 3. Contracts

- 模型配置：示例顶部常量只保存占位 Key；真实 Key 不提交。
- 基础设施配置：使用 README 明确的环境变量，例如 `VECTOR_BASE_URL`、`MYSQL_DSN`、`RELEASE_CONTROL_BASE_URL`。
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

### 6. Tests Required

- 仓库级 `catalog_test.go` 验证目录、TODO 连续性和真实构造入口符号存在。
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
