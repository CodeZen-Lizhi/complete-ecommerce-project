项目demo
 公共的封装的增删改查方法
规范使用方式  按层级调用过去 
返回 controller  写一个 result

事务

log/slog
错误嵌套 return fmt.Errorf("failed to read config file: %v", err)

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


//检查是否实现接口 的断言
// 定义一个接口
type IAnimal interface {
Speak() string
}

// 定义一个结构体
type Dog struct {
Name string
}

// Dog 实现 IAnimal 接口
func (d *Dog) Speak() string {
return "Woof!"
}

// 编译时检查，确保 *Dog 实现了 IAnimal 接口
var _ IAnimal = (*Dog)(nil)

