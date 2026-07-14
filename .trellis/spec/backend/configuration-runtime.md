# 配置与运行时规范

## 启动与关闭顺序

`main.go` 是唯一生产进程入口，当前顺序为：

```text
config.Init
-> logger.InitLogConfig
-> util.Init（Snowflake）
-> redis.Init
-> mysql.InitMySQL
-> container.GetInstance
-> router.InitTotalRouter
-> HTTP Server
```

Redis 和 MySQL 在初始化成功后立即注册 `defer Close`。HTTP Server 使用 `signal.NotifyContext` 监听退出信号，并用 10 秒超时执行 `Shutdown`。新增长期资源必须遵循“初始化失败立即返回、初始化成功注册关闭”的生命周期。

## 配置加载

- `APP_ENV` 决定读取 `configs/config.<env>.yaml`，缺省为 `dev`。
- 配置结构和 `mapstructure` 标签定义在 `internal/config/config.go`。
- Viper 使用 `APP` 环境变量前缀；Redis 地址另有显式 `REDIS_HOST` 绑定。
- `UnmarshalConfig` 每次创建新的 `Config` 并执行 `validateConfig`。
- 新配置项必须同时更新结构体、YAML 示例、校验逻辑和实际消费者。

Viper 会监听文件变更，但 logger、Redis、MySQL 等组件在启动时读取配置并通过 `sync.Once` 初始化。不要宣称这些资源会随 YAML 热更新；需要动态生效时必须单独设计重建或原子切换机制。

## MySQL 与 Redis 生命周期

- `internal/mysql.InitMySQL` 使用配置创建 GORM，并设置最大连接数、空闲连接数和连接寿命。
- `internal/mysql.DB` 只在初始化成功后赋值；业务层通过 Container/Repository 使用它。
- `internal/redis.Init` 创建单例 Client，并用 5 秒 context 执行 Ping。
- `redis.Client()` 在未初始化时会 panic，只能在 `Init` 成功后调用。
- 外部调用必须携带调用链 ctx 或明确超时，不能无期限阻塞。

## HTTP 进程设置

- 服务端口来自 `config.Cfg.App.Port`。
- `http.Server` 当前配置 `ReadHeaderTimeout: 5s`。
- 非 `prod` 环境会在 `localhost:6060` 启动 pprof；不得暴露到公网地址。
- `http.ErrServerClosed` 是正常停机信号，不记录为启动失败。

## 密钥与环境配置

- `configs/config.dev.yaml` 是本地开发配置，不是生产 Secret 存储。
- 当前 JWT 工具会拒绝默认弱密钥，并要求至少 32 字符，参考 `util/jwtutil.go`。
- 不新增或提交真实数据库密码、JWT Secret、API Key、Token 和证书。
- 需要本地敏感值时使用未跟踪的环境变量或本地覆盖文件；日志与错误消息不得打印它们。

## 配置变更检查

1. 搜索配置键在 `Config`、YAML、环境绑定、校验和消费者中的全部引用。
2. 确认缺省值、零值和非法值的行为。
3. 确认热更新是否真的能被消费者观察到。
4. 运行配置加载测试或最小启动烟测；涉及 MySQL/Redis 时记录依赖服务前提。
