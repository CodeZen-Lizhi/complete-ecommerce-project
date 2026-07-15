# 跨层思考指南

> 用于修改前确定影响面。具体编码规则以 `../backend/index.md` 及其子指南为准。

## 指南索引

| 指南 | 何时使用 |
| --- | --- |
| [跨层变更指南](./cross-layer-thinking-guide.md) | 新增业务模块，修改 API/DTO、认证、事务、配置或运行时链路 |
| [代码复用指南](./code-reuse-thinking-guide.md) | 准备新增 helper、封装、全局实例、响应格式或重复逻辑 |

## 快速触发条件

出现以下任一情况，先读跨层指南：

- 改动覆盖 Router、Handler、Service、Repository 中两个以上层。
- 修改 `UserService`、`UserRepository`、`Container`、`Result` 或配置结构。
- 新增路由组、中间件、权限规则、事务、缓存或外部模型调用。
- 修改数据库字段、索引、软删除、ID 策略或响应状态码。

出现以下任一情况，先读复用指南：

- 准备新增 `util`、`helper`、`common` 或全局单例。
- 发现相似 JSON 响应、分页、错误包装、Context Key 或连接初始化。
- 想在 Handler 直接访问 GORM/Redis，或在 Service 直接生成 HTTP 响应。
- 同一配置或常量需要在多个文件重复声明。

## 证据规则

所有影响面结论以源码搜索、调用方、配置、数据库事实和实际测试为准。当前生产业务层仍缺少系统化测试，仓库也没有迁移系统；缺少证据时应明确记录缺口，不把 AI 示例测试、README 或 Demo Handler 当成生产行为覆盖。
