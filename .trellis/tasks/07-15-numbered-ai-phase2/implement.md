# 阶段 2：LLM API 与 Prompt 编号练习执行计划

## 实现清单

- [x] 生成 `01_basic_chat`：基础模型调用。
- [x] 生成 `02_eino_generate`：Eino ChatModel Generate。
- [x] 生成 `03_prompt_roles`：System 角色与 User Prompt。
- [x] 生成 `04_in_memory_multi_turn`：程序内多轮对话。
- [x] 生成 `05_redis_session_history`：Redis 会话历史。
- [x] 生成 `06_session_history_limits`：历史截断、TTL 与用户隔离。
- [x] 生成 `07_chat_template`：Eino ChatTemplate 与变量。
- [x] 生成 `08_streaming_chat`：ChatModel Stream、EOF 与取消。
- [x] 生成 `09_structured_json_output`：固定 JSON、解析与字段校验。
- [x] 生成 `10_call_governance`：超时、重试、限流与调用统计。
- [x] 生成 `11_model_provider_adapter`：自定义接口与模型配置切换。

## 一致性检查

- [x] 验证目录编号从 `01` 连续到 `11`。
- [x] 验证每个目录包含 README，章节和 TODO 编号连续。
- [x] 验证默认配置不会连接外部服务或产生付费请求。
- [x] 搜索本阶段旧路径、重复目录、真实密钥和未处理错误。

## 验证

```bash
gofmt -w <阶段 2 新增或迁移的 Go 文件>
go test -timeout=60s ./examples/ai/phase2/...
go vet ./examples/ai/phase2/...
git diff --check
```

阶段 2 完成后执行 `trellis-check` 和 `go-review`；涉及 SQL/Schema/查询时追加 `sql-code-review`。
