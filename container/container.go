package container

import (
	"ecommerce/repository"
	"ecommerce/service"
	"errors"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

// GetInstance 获取容器单例（显式注入 DB）
func GetInstance(db *gorm.DB) *Container {
	once.Do(func() {
		// 手动组装依赖（依赖关系清晰）
		userRepository := repository.NewUserRepository(db)
		userService := service.NewUserService(userRepository)

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
