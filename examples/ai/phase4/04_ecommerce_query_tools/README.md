# 阶段 4 练习 4：电商业务查询工具

## 练习目标

将商品、库存、订单和物流查询包装为最小只读工具接口。

## 前置知识

- 模型只生成 Tool Call，应用负责校验、授权和执行。
- 模型参数是不可信输入，Context 中的身份才是可信边界。
- 写操作需要确认、幂等、审计和明确失败处理。

## TODO 顺序

核心接口与函数签名位于 `exercise.go`，`main.go` 只保留安全启动入口。

1. TODO 1：为每类查询定义窄接口和只读 DTO。
2. TODO 2：校验商品/订单 ID 与分页上限。
3. TODO 3：注入 Fake Service，不在 Tool 中直接访问 GORM。
4. TODO 4：根据可信身份过滤订单和物流数据。
5. TODO 5：限制返回字段和列表大小，并覆盖 NotFound 与依赖错误。

## 开始练习

```bash
go run ./examples/ai/phase4/04_ecommerce_query_tools
```

骨架默认使用本地接口和安全配置，不连接生产数据库，也不执行真实写操作。

## 验证方式

```bash
gofmt -w examples/ai/phase4/04_ecommerce_query_tools/*.go
go test -timeout=60s ./examples/ai/phase4/04_ecommerce_query_tools
go vet ./examples/ai/phase4/04_ecommerce_query_tools
```

## 完成标准

- 参数、身份、租户、权限和结果上限都在工具边界验证。
- Context 取消和依赖错误清楚传播，错误不泄露敏感内部信息。
- 只读调用可测试，写操作默认禁用并具备确认、幂等和审计。
- 不让模型自由决定可信身份或绕过固定业务规则。

## 暂不实现

- 真实生产数据库接入和写操作。
