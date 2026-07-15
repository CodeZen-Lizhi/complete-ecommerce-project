# 目录与分层

## 实际目录

```text
main.go                 应用启动、资源初始化、HTTP Server 与优雅停机
configs/                按环境区分的 YAML 配置
container/              手动组装 Repository 和 Service
router/                 路由分组与中间件挂载
middleware/             Gin 横切能力：请求 ID、日志、恢复、认证
handler/                HTTP 输入解析、轻量校验、响应映射
service/                业务规则、密码处理、事务边界
repository/             GORM 数据访问
model/                  持久化模型和现有传输结构
internal/config/        Viper 加载、反序列化和校验
internal/logger/        slog 初始化
internal/mysql/         GORM 连接和事务入口
internal/redis/         go-redis 客户端生命周期
internal/response/      统一 API 响应实现
util/                   JWT、雪花 ID 和跨包常量
examples/               与生产启动链隔离的学习示例
docs/                   项目文档和学习路线
```

证据：`README.md` 的项目结构、`main.go` 的初始化顺序和 `go list ./...` 的实际包集合。

## 依赖方向

- `router` 只注册路由和中间件，参考 `router/total_router.go`、`router/user_router.go`。
- `handler` 依赖 Gin、Container、Service 契约和响应辅助函数，参考 `handler/user_handler.go`。
- `service` 依赖 Repository 接口和领域模型，业务事务可调用 `internal/mysql.Transaction`。
- `repository` 依赖 GORM 和 `model`，不依赖 Handler、Router 或 Gin。
- `container/container.go` 是生产依赖组装点：`DB -> UserRepository -> UserService`。
- `main.go` 只负责进程生命周期，不放业务路由或查询逻辑。

禁止形成反向依赖，例如 Repository 导入 Service、Service 写入 Gin 响应、Model 读取 HTTP Header。

## 新增业务模块

当新模块需要完整持久化链路时，按以下顺序落位：

1. 在 `model/<name>.go` 定义数据库模型，并确认真实表名和字段。
2. 在 `repository/<name>_repository.go` 定义数据访问接口、实现和构造函数。
3. 在 `service/<name>_service.go` 定义业务接口、实现和构造函数。
4. 在 `container/container.go` 增加 Service 字段并完成手动组装。
5. 在 `handler/<name>_handler.go` 实现协议适配。
6. 在 `router/<name>_router.go` 注册 public/private/admin 路由。
7. 为受影响层补测试，并验证完整请求路径。

参考完整模块：User 模块的 `model/user.go`、`repository/user_repository.go`、`service/user_service.go`、`handler/user_handler.go` 和 `router/user_router.go`。

## 命名与文件边界

- 包名使用小写单词；当前模块导入前缀为 `ecommerce/...`。
- 分层文件沿用 `<feature>_handler.go`、`<feature>_service.go`、`<feature>_repository.go`、`<feature>_router.go`。
- 导出类型和函数使用 Go 的 PascalCase；现有接口采用 `UserService` / `UserRepository`，实现采用 `UserServiceImpl` / `UserRepositoryImpl`。
- 每个新增函数和方法无论是否导出，都写简洁中文职责注释；详细格式遵循 `quality-guidelines.md` 的“函数与方法注释约定”。
- 仅供包内路由注册的函数保持小写，例如 `registerUserPublicRoutes`。

## 独立学习示例

`examples/ai/phase2/01_basic_chat` 是学习代码，不接入 `main.go`、Container 或生产 Router。后续学习练习遵循 `docs/ai-learning/go-eino-ai-engineer-roadmap.md` 的“骨架 + 中文 TODO”规范；只有明确要求产品化时才迁入正式分层。

### AI 练习目录约定

- 学习路线使用 `examples/ai/phase1/` 至 `phase7/`，阶段内目录固定为两位数字前缀的 `NN_<slug>`，例如 `07_chat_template`。
- 阶段 1 目录只做现有电商项目的任务导航，不复制生产业务代码；阶段 2–7 提供隔离的可编译 Go 骨架。
- 每个练习必须包含中文 `README.md`；阶段 2–7 还必须包含 `main.go`，阶段 3–7 使用 `exercise.go` 保存题目专属接口和函数签名。核心步骤使用连续编号的中文 TODO，并在未完成时明确退出。
- 默认占位配置必须在模型、Redis、数据库或向量库客户端创建前校验，避免骨架运行产生真实请求或外部写入。
- 新增、删除或重排练习后必须更新 `examples/ai/README.md`，并运行 `go test -timeout=60s ./examples/ai` 验证 7 个阶段、60 个目录、连续编号和必需文件。

正确目录：

```text
examples/ai/phase3/01_embedding_similarity/
├── README.md
└── main.go
```

错误目录：

```text
examples/ai/phase3/embedding/     # 缺少阶段内顺序
examples/ai/common/               # 让练习之间形成隐式答案依赖
```

## 不要复制的现有占位

- `handler/product_handler.go` 中硬编码商品、分类和分页结果。
- `handler/user_handler.go` 中验证码、订单、更新资料和修改密码的演示响应。
- Handler 内忽略 `strconv` 错误或直接使用 `PostForm` 代替明确 DTO 校验的写法。

这些代码描述当前 Demo 状态，不代表新增功能的推荐结构。
