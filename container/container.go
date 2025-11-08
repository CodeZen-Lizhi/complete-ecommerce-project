package container

import (
	"ecommerce/repository"
	"ecommerce/service"
	"errors"
	"github.com/gin-gonic/gin"
	"sync"
)

// Container 依赖容器
type Container struct {
	UserService service.UserService
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
		userRepository := repository.NewUserRepository()
		userService := service.NewUserService(userRepository)

		// 5. 组装容器
		instance = &Container{
			UserService: userService,
		}
	})
	return instance
}

// GetContainer 从 Context 中获取容器
func GetContainer(c *gin.Context) (*Container, error) {
	val, exists := c.Get(ContainerKey)
	if !exists {
		return nil, errors.New("容器未注入")
	}

	ctn, ok := val.(*Container)
	if !ok {
		return nil, errors.New("容器类型错误")
	}
	return ctn, nil
}
