# 质量规范

## 当前基线

- Go 版本：`go.mod` 中的 `go 1.25`。
- 格式化：标准 `gofmt`。
- 基础验证：`go test ./...` 和 `go vet ./...`。
- 仓库当前没有 `_test.go`，因此测试覆盖是明确缺口，不得声称已有单元测试保护。
- 没有仓库级 golangci-lint、Makefile 或已跟踪 Dockerfile；不要编造对应命令。

## 实现要求

- 保持最小、聚焦的 diff，不做无关格式化或重构。
- 优先使用现有 Container、响应辅助函数、事务封装、Logger 和 Context Key。
- 导出符号写简洁中文注释；注释解释职责、约束或取舍。
- 函数保持单一职责、浅嵌套和明确错误返回。
- 新依赖前先确认标准库和现有依赖不能满足需求，并说明维护与安全影响。
- 不修改或覆盖用户未提交的工作区改动。

## 测试策略

按风险为新增或修改的行为补最接近边界的测试：

- Service：注入 Fake/Mock Repository，覆盖成功、业务错误、依赖错误和事务回滚。
- Handler/中间件：使用 `httptest` 和 Gin，覆盖输入校验、响应结构、状态码、认证和 Abort 行为。
- Repository：使用隔离数据库或受控 GORM 测试环境，覆盖 NotFound、软删除、更新零值和事务。
- JWT/工具：表驱动测试覆盖合法、过期、签名算法错误、issuer/audience 错误和空输入。
- 并发、共享状态或 goroutine 变更追加 `go test -race ./...`。

后端单测应设置不超过 60 秒的合理超时。测试必须能在删除被测实现后失败，避免只复述实现的同义断言。

## 安全红线

- 不硬编码或提交真实 Secret、Token、API Key、证书和生产数据库凭证。
- SQL 和命令调用参数化；动态标识符使用白名单。
- 权限必须在后端根据可信身份执行，不能信任客户端 `Role` Header。
- 外部输入在 Handler/边界校验；输出 DTO 不暴露密码等敏感字段。
- 日志不包含 Authorization、Cookie、密码、Token 或完整敏感数据。

## 已知 Demo/技术债处理

以下现有代码不能作为新实现范例：

- `handler/product_handler.go` 的硬编码商品/分类与忽略转换错误。
- `handler/user_handler.go` 的假验证码、假订单和未实现密码修改。
- `AdminAuthMiddleware` 基于请求 Header 的管理员判断。
- Service 到 Repository 的 ctx 传播缺失。
- Login 历史密码升级吞掉 Update 错误。

除非任务明确要求，不顺手大范围修复这些问题；但新增代码不得扩散同类模式。

## 修改后验证

最小命令集：

```bash
gofmt -w <本次修改的 Go 文件>
go test ./...
go vet ./...
```

按改动追加：

- HTTP：真实 Router 的 `httptest` 或最小接口烟测。
- 数据库：事务/查询测试和必要的 `EXPLAIN`。
- 并发：`go test -race ./...`。
- 配置/启动：使用明确环境的最小启动与优雅停机检查。
- AI 学习示例：确保骨架可编译，未完成时明确提示且不产生真实付费调用。

## Review 清单

1. Correctness：成功、空值、边界、失败和取消路径是否完整。
2. Layering：Handler、Service、Repository 的职责是否越界。
3. Data/Security：事务、软删除、权限、敏感数据和参数化是否正确。
4. Concurrency/Performance：是否有泄漏、竞态、N+1、无界分页或循环远程调用。
5. Compatibility：HTTP/业务码、DTO、配置、表字段和调用方是否兼容。
6. Verification：测试、vet、构建或烟测结果是否真实执行并说明覆盖不足。

Go 修改交付前按根目录 `AGENTS.md` 使用 `go-review`；涉及 SQL/Schema/查询/索引时再追加 SQL 专项 review。
