# 阶段 4：Tool Calling 编号练习执行计划

## 实现清单

- [x] 生成 `01_local_readonly_tool`：本地只读工具与 ToolsNode。
- [x] 生成 `02_tool_argument_validation`：必填、类型、范围与未知字段校验。
- [x] 生成 `03_context_identity_authorization`：Context 身份、租户与再次授权。
- [x] 生成 `04_ecommerce_query_tools`：商品、库存、订单与物流查询工具。
- [x] 生成 `05_multi_tool_orchestration`：多工具注册与连续调用。
- [x] 生成 `06_tool_resilience`：超时、下游失败、返回上限与有限重试。
- [x] 生成 `07_safe_write_tool`：取消订单/退款模拟、确认、幂等与审计。

## 一致性检查

- [x] 验证目录编号从 `01` 连续到 `07`。
- [x] 验证每个目录包含 README，章节和 TODO 编号连续。
- [x] 验证默认配置不会连接外部服务或产生付费请求。
- [x] 搜索本阶段旧路径、重复目录、真实密钥和未处理错误。

## 验证

```bash
gofmt -w <阶段 4 新增或迁移的 Go 文件>
go test -timeout=60s ./examples/ai/phase4/...
go vet ./examples/ai/phase4/...
git diff --check
```

阶段 4 完成后执行 `trellis-check` 和 `go-review`；涉及 SQL/Schema/查询时追加 `sql-code-review`。
