# 阶段 2 练习 11：模型调用治理

## 练习目标

在统一入口增加 Context 超时、有限重试、限流和调用统计。

## TODO 顺序

1. TODO 1：在外层限流器获得许可，等待时响应取消。
2. TODO 2：为每次尝试创建并及时释放独立超时 Context。
3. TODO 3：分类错误，只重试明确可恢复错误。
4. TODO 4：执行有上限的指数退避，并在等待时响应取消。
5. TODO 5：记录尝试次数、总耗时、Token 和成功状态，日志不含敏感内容。
6. TODO 6：成功返回响应和指标；耗尽重试时保留最后错误链。

## 开始练习

```bash
go run ./examples/ai/phase2/11_call_governance
```

默认 API Key 或外部服务配置保持占位符时，程序应明确提示配置未完成，不得连接外部服务或产生付费请求。

## 验证方式

```bash
gofmt -w examples/ai/phase2/11_call_governance/*.go
go test -timeout=60s ./examples/ai/phase2/11_call_governance
go vet ./examples/ai/phase2/11_call_governance
```

## 完成标准

- TODO 编号连续，错误使用 `%w` 保留底层原因。
- Context 超时和取消能够传播到模型或网络边界。
- 不打印 API Key、完整敏感 Prompt 或 Authorization。
- 删除核心实现后，测试或默认运行能够清楚暴露练习未完成。

## 暂不实现

- 分布式限流、熔断平台和复杂成本结算。
