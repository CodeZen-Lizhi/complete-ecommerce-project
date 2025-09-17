package user

import (
	"ecommerce/model"
	"ecommerce/repository/user"
)

// 接口
type UserService interface {
	FindByID(id uint) (*model.User, error)
	izExist(username string) (bool, error)
}

// 接口实现类
type userServiceImpl struct {
	userRepository user.UserRepository
}

// 确保实现类实现了接口
var _ UserService = (*userServiceImpl)(nil)

func NewUserService(userRepository user.UserRepository) UserService {
	return &userServiceImpl{
		userRepository: userRepository,
	}
}

func (u *userServiceImpl) FindByID(id uint) (*model.User, error) {
	return u.userRepository.FindByID(id)
}

func (u *userServiceImpl) izExist(username string) (bool, error) {
	exist, err := u.userRepository.izExist(username)
	return exist, err
}
