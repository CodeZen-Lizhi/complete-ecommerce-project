package container

import (
	UserDao "ecommerce/repository/user"
	UserService "ecommerce/service/user"
	"sync"
)

// Container 依赖容器
type Container struct {
	UserService *UserService.UserServiceImpl
}

var (
	instance *Container
	once     sync.Once // 保证单例初始化
)

// ContainerKey 用于 Gin Context 存储的键
const ContainerKey = "container"

// GetInstance 获取容器单例
func GetInstance() *Container {
	once.Do(func() {
		userRepository := UserDao.NewUserRepository()
		userService := UserService.NewUserService(userRepository)

		// 5. 组装容器
		instance = &Container{
			UserService: userService,
		}
	})
	return instance
}
