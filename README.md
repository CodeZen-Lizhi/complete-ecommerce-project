项目demo
todo 
公共的封装的增删改查方法
规范使用方式  按层级调用过去 
返回 controller  写一个 result
事务
错误嵌套 return fmt.Errorf("failed to read config file: %v", err)
先手写注入，后面在用框架

	// 最底层的错误
	rootErr := errors.New("网络连接超时")

	// 第一层包装
	wrappedErr1 := fmt.Errorf("数据库连接失败: %w", rootErr)

	// 第二层包装
	wrappedErr2 := fmt.Errorf("用户服务初始化失败: %w", wrappedErr1)

	// 第三层包装
	finalErr := fmt.Errorf("应用启动失败: %w", wrappedErr2)

	fmt.Println("最终错误:", finalErr)
	// 输出: 应用启动失败: 用户服务初始化失败: 数据库连接失败: 网络连接超时

	// 使用 errors.Unwrap() 解包错误
	fmt.Println("\n解包一层:", errors.Unwrap(finalErr))
	fmt.Println("解包两层:", errors.Unwrap(errors.Unwrap(finalErr)))
	fmt.Println("解包三层:", errors.Unwrap(errors.Unwrap(errors.Unwrap(finalErr))))

	// 使用 errors.Is() 检查错误链中是否包含特定错误
	if errors.Is(finalErr, rootErr) {
		fmt.Println("\n错误链中包含'网络连接超时'错误")
	}
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

