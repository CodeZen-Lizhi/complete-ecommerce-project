# 阶段 2 练习 3：System Prompt 角色对照

## 这次练什么

前两个练习已经完成了基础 HTTP 模型调用和 Eino `Generate`。这次只研究一个变量：**System Prompt 如何改变同一个 User Prompt 的回答方式**。

你会让模型分别扮演：

1. 面向 Go 初学者的导师。
2. 面向资深 Go 工程师的技术顾问。

两次调用必须使用完全相同的 User Prompt、模型和调用方式。否则无法判断输出变化究竟来自 System Prompt，还是来自其他变量。

## 数据流

```text
相同 User Prompt
├── System Prompt：初学者导师 → Generate → 回答 A
└── System Prompt：资深工程师顾问 → Generate → 回答 B

比较回答 A 与回答 B
```

## 准备

当前示例使用本地 OpenAI-compatible 服务的 Responses API：

- Base URL：`http://localhost:8084/v1`
- Endpoint：`POST /v1/responses`
- 模型：`gpt-5.4-mini`
- API Key：使用本地 AI 网关的自定义 Key，直接配置在本练习源码中

该 Key 只用于用户本机部署的 `localhost:8084` 网关，不是公网服务凭据。

本练习使用 Eino 官方 `agenticopenai.NewResponsesModel`。System 和 User Prompt
分别通过 `schema.SystemAgenticMessage`、`schema.UserAgenticMessage` 构造，模型
回答从 `AgenticMessage.ContentBlocks` 中的 `AssistantGenText` 提取。

## 开始练习

```bash
go run ./examples/ai/phase2/prompt_roles
```

程序会按顺序请求两个场景，并输出场景名称、System Prompt 和模型回答。

## 观察重点

完成后比较两组回答：

- 使用的术语是否不同？
- 代码示例的复杂度是否不同？
- 是否解释了 interface 的使用成本？
- 是否针对目标受众调整了篇幅和表达方式？
- System Prompt 是否真的约束了回答，还是只增加了无关措辞？

## 本练习暂时不做

- 不保存 Assistant 历史。
- 不实现多轮对话。
- 不接入 Redis。
- 不使用 ChatTemplate 或 Stream。
- 不实现重试、限流和 Token 统计。

这些内容会在阶段 2 的后续练习中逐步加入。

## 验证

```bash
gofmt -w examples/ai/phase2/prompt_roles/main.go
go test -timeout=60s ./...
go vet ./...
```
