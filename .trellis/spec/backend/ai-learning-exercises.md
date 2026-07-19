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

- 模型配置：需要真实模型调用的练习在 `exercise.go` 或 `main.go` 顶部提供 `BaseURL`、`API Key` 和模型名常量，学习者直接修改这些入口。
- 指向 `localhost` 等本地 OpenAI-compatible 服务、且经项目所有者明确确认不敏感的本地专用 Key，可以作为练习常量提交；不得仅凭 `sk-` 等字符串外形把它判定为外部真实密钥或擅自恢复成占位符。
- 外部、共享或生产模型服务的真实 Key 仍不得提交，必须保留占位符或使用环境变量。无法从地址和项目约定判断 Key 性质时，提交前向项目所有者确认。
- 基础设施配置：Qdrant、BM25、MySQL、Reranker、发布控制等练习也在顶部提供服务专属占位常量，README 必须列出应修改的常量名。
- Secret 管理与 Docker Compose 练习是例外：它们的学习目标就是环境变量和 Secret 生命周期，不得把真实密钥改成源码常量。
- README 必须包含真实运行命令、是否产生模型费用、预期观察项和失败场景。
- 用户明确要求时，练习保持 TODO 骨架，不直接交付核心答案。

### 4. Validation & Error Matrix

| 条件 | 行为 |
| --- | --- |
| Context 为 nil | 返回明确错误 |
| API Key 为空或仍为占位值 | 真实调用前失败 |
| 本地专用 Key 已由项目所有者确认不敏感 | 允许从顶部常量读取并提交 |
| Key 属于外部、共享或生产服务 | 不得提交，改用占位符或环境变量 |
| 外部服务地址为空 | 构造 Client 前失败 |
| SDK 构造或真实调用失败 | 使用 `%w` 保留底层错误链 |
| 只实现 Fake/Mock/等价接口 | 不满足 README 完成标准 |

### 5. Good/Base/Bad Cases

- Good：骨架保留小接口用于理解边界，同时提供真实 SDK 构造函数 TODO，并要求运行真实组件闭环。
- Good：`BaseURL` 指向本地服务，项目所有者已确认对应 Key 仅用于本地练习，因此配置保留在顶部常量并随练习提交。
- Base：目录、README、Go 骨架可编译，未完成时明确返回 `errExerciseIncomplete`。
- Bad：README 写“Eino 或等价实现”，最终只调用本地 Map、Fake 模型或手写状态机。
- Bad：多个编号 TODO 连续堆在入口函数里，但对应类型和函数没有可填写位置。
- Bad：未区分本地专用凭证和外部真实凭证，仅根据 Key 的字符串格式阻止提交或擅自覆盖用户配置。

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

## Scenario: Structured Outputs 约束分层

### 1. Scope / Trigger

- 适用于 `examples/ai/` 中使用 OpenAI-compatible `response_format.type=json_schema` 的结构化输出练习。
- 目标是区分模型原始 JSON、解析后的 Go 结构体和终端展示文本，避免把展示格式误判为 Schema 未生效。

### 2. Signatures

- `buildResponseFormat(mode structuredOutputMode) (*openai.ChatCompletionResponseFormat, error)` 负责构造模型请求的响应格式，不负责打印或校验自然语言质量。
- `decodeAndValidate(raw []byte) (structuredAnswer, error)` 负责严格解析模型原始 JSON，并执行 Schema 之外的业务校验。
- `main()` 可以按人类可读格式展示解析结果；需要验证原始 JSON 时必须在模型边界检查 `response.Content`，不得用展示文本代替。

### 3. Contracts

- JSON Schema 负责可机器判定的结构约束：字段、类型、必填项、额外字段、长度、数量、数值范围、枚举和支持的组合条件。
- System Prompt 和字段 `description` 负责语言、语气、简洁程度、禁止 Markdown 等生成要求；原生模式无需重复 JSON 字段示例，但仍应保留这些语义要求。
- 应用端校验负责必须严格保证但 Schema 或服务端不能可靠表达的规则，例如禁止 Markdown 代码围栏、数组元素不得为空、只允许一个 JSON 值。
- `Strict: true` 只表示严格遵守服务端接受的 Schema，不代表所有自然语言要求都会成为硬约束。

### 4. Validation & Error Matrix

| 条件 | 行为 |
| --- | --- |
| 缺少必填字段、字段类型错误或出现未知字段 | 返回 JSON 结构错误 |
| 出现第二个 JSON 值或尾随非空内容 | 返回 JSON 结构错误 |
| 字符串为空、数组为空或数值越界 | 返回业务字段错误 |
| 字符串包含明确禁止的 Markdown 代码围栏 | 应用端校验返回业务字段错误 |
| 模型服务不支持 `json_schema` | 保留创建或调用错误，不静默降级到 Prompt 模式 |
| 终端使用 `fmt.Printf` 展示已解析字段 | 视为展示格式，不据此判断原始 JSON 是否符合 Schema |

### 5. Good/Base/Bad Cases

- Good：Schema 限制字段、类型和范围；Prompt 描述纯文本要求；应用端再次拒绝代码围栏。
- Base：模型返回单个合法 JSON 对象，能够解码为 `structuredAnswer`。
- Bad：只写 `description=禁止Markdown`，却把它当作服务端强制校验规则。
- Bad：看到终端输出“摘要/关键点/置信度”后，误认为模型没有返回 JSON，而不检查模型边界和解析结果。

### 6. Tests Required

- `buildResponseFormat` 测试必须断言原生模式包含 `json_schema`、`Strict: true`、`required` 和 `additionalProperties:false`。
- `decodeAndValidate` 测试必须覆盖未知字段、多个 JSON、缺失字段、空数组、空数组元素和数值越界。
- 新增禁止 Markdown 等硬业务规则时，必须增加对应非法输入测试；仅检查 Schema 中存在 `description` 不算完成。
- 真实服务烟测必须区分原始响应和 `main()` 展示输出，并记录服务端是否支持所用 Schema 关键字。

### 7. Wrong vs Correct

#### Wrong

```go
type answer struct {
	Summary string `json:"summary" jsonschema:"description=禁止Markdown"`
}
```

仅靠 `description` 不能保证模型输出不包含 Markdown。

#### Correct

```go
type answer struct {
	Summary string `json:"summary" jsonschema:"minLength=1,maxLength=200,description=使用简洁纯文本，禁止Markdown"`
}

if strings.Contains(payload.Summary, "```") {
	return answer{}, fmt.Errorf("%w: summary 禁止包含 Markdown 代码围栏", errBusinessField)
}
```

Schema 提供可机器判断的边界，Prompt/`description` 指导生成，应用端校验硬性禁止项。
