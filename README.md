# Go 电商项目模板

一个基于 Go 的电商后端项目模板，采用分层架构设计，适合作为新项目的起点。

## 项目特点

- ✅ 清晰的分层架构（Handler → Service → Repository）
- ✅ 手动依赖注入（简单明了，易于理解）
- ✅ 完整的用户模块示例（注册、登录、认证）
- ✅ 事务支持（GORM 事务封装）
- ✅ 统一的错误处理和响应格式
- ✅ 结构化日志（slog）
- ✅ 中间件支持（CORS、认证、日志、恢复）
- ✅ 路由分组（公共/私有/管理员）

## 技术栈

| 技术类别 | 使用组件     | 说明                      |
|----------|-------------|---------------------------|
| 配置读取 | viper       | 用于读取和管理应用配置文件 |
| 日志     | log/slog    | Go 标准库中的结构化日志工具 |
| Web框架  | Gin         | 高性能的 HTTP Web 框架     |
| ORM      | GORM        | 强大的数据库 ORM 库        |
| 热重载   | Air         | Go 应用开发时的实时重载工具 |
| Redis    | go-redis    | Redis 数据库的 Go 客户端   |
| 依赖注入 | 手动注入     | 简单直接，无需额外框架     |

## 项目结构

```
.
├── configs/              # 配置文件
├── container/            # 依赖注入容器
├── handler/              # HTTP 处理层
├── service/              # 业务逻辑层
├── repository/           # 数据访问层
├── model/                # 数据模型
├── router/               # 路由定义
├── middleware/           # 中间件
├── internal/             # 内部包
│   ├── config/          # 配置管理
│   ├── logger/          # 日志管理
│   ├── mysql/           # MySQL 连接
│   └── redis/           # Redis 连接
└── util/                 # 工具类

```

## 依赖注入

本项目采用**手动依赖注入**，简单明了，易于理解和维护。

### 依赖关系

```
DB → UserRepository → UserService → Handler
```

### 为什么用手动注入？

- ✅ 代码清晰，依赖关系一目了然
- ✅ 无需学习额外框架（如 Wire、Fx）
- ✅ 调试方便，无生成代码
- ✅ 符合 Go 社区主流实践（60% 项目采用）
- ✅ 适合中小型项目（< 20 个依赖）

### 示例代码

```go
// container/container.go
func GetInstance(db *gorm.DB) *Container {
    // 手动组装依赖
    userRepository := repository.NewUserRepository(db)
    userService := service.NewUserService(userRepository)

    return &Container{
        UserService: userService,
    }
}

// main.go
mysql.InitMySQL()
ctn := container.GetInstance(mysql.DB)  // 显式传入依赖
```

### 如果需要 Wire

当依赖数量 > 50 个时，可以考虑引入 Google Wire：
- 参考：https://github.com/google/wire
- 示例：go-kratos 框架

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 配置数据库

编辑 `configs/config.dev.yaml`：

```yaml
mysql:
  dsn: "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
```

### 3. 运行项目

```bash
# 开发模式（热重载）
air

# 或直接运行
go run main.go
```

## 架构说明

### 分层架构

```
Handler (HTTP 层)
   ↓
Service (业务逻辑层)
   ↓
Repository (数据访问层)
   ↓
Database
```

### 模块示例

项目包含完整的 User 模块作为示例：

- **Handler**: `handler/user_handler.go` - HTTP 请求处理
- **Service**: `service/user_service.go` - 业务逻辑（注册、登录、密码加密）
- **Repository**: `repository/user_repository.go` - 数据库操作
- **Model**: `model/user.go` - 数据模型

### 事务支持

```go
// service/user_service.go
func (u *UserServiceImpl) Register(ctx context.Context, user *model.User) error {
    return mysql.Transaction(ctx, func(tx *gorm.DB) error {
        // 使用事务 DB
        repo := u.userRepository.WithDB(tx)

        // 检查用户是否存在
        exist, err := repo.IzExist(user.Username)
        if err != nil {
            return err
        }
        if exist {
            return ErrUserAlreadyExists
        }

        // 创建用户
        return repo.Create(user)
    })
}
```

## 路由分组

```go
// 公共路由（无需登录）
/api/public/user/register
/api/public/user/login
/api/public/products

// 私有路由（需登录）
/api/private/user/info
/api/private/orders

// 管理员路由（需管理员权限）
/api/admin/products
/api/admin/orders
```

## 开发建议

### 添加新模块

1. 创建 Model（`model/xxx.go`）
2. 创建 Repository（`repository/xxx_repository.go`）
3. 创建 Service（`service/xxx_service.go`）
4. 创建 Handler（`handler/xxx_handler.go`）
5. 注册路由（`router/xxx_router.go`）
6. 在 Container 中注入依赖

### 测试建议

- Repository 层：使用接口，便于 Mock
- Service 层：注入 Mock Repository 进行单元测试
- Handler 层：使用 `httptest` 进行集成测试

## TODO

- [ ] 公共的封装的增删改查方法
- [ ] 更多模块示例
- [ ] 单元测试示例
- [ ] API 文档（Swagger）
