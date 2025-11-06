项目demo
todo 
公共的封装的增删改查方法
规范使用方式  按层级调用过去 
返回 controller  写一个 result
事务
GinRecovery 这还是不完善，没有堆栈信息
先手写注入，后面在用框架
jwt-go 引入
Makefile 核心是通过 make 命令执行预定义的指令，简化 Go 项目的编译、打包、测试等操作。
 试试 Makefile
Google Wire 引入
试试 air
技术栈：

| 技术类别  | 使用组件     | 说明                  |
|-------|----------|---------------------|
| 配置读取  | viper    | 用于读取和管理应用配置文件       |
| 日志    | log/slog | Go 标准库中的结构化日志记录工具   |
| Web框架 | GIN      | 高性能的 HTTP Web 框架    |
| ORM   | GORM     | 强大的数据库 ORM 库        |
| 热重载   | AIR      | Go 应用开发时的实时重载工具     |
| Redis | Go-Redis | Redis 数据库的 Go 语言客户端 |
| 依赖注入  | Google Wire | 自动依赖注入框架 |

