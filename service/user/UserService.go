package user

import (
	"ecommerce/model"
	"ecommerce/repository/user"
)

// 接口
type UserService interface {
	FindByID(id uint64) (*model.User, error)
	IzExist(username string) (bool, error)
	Create(user *model.User) error
	Login(username string, password string) (*model.User, error)
}

// 接口实现类
type UserServiceImpl struct {
	userRepository user.UserRepository
}

// 确保实现类实现了接口
var _ UserService = (*UserServiceImpl)(nil)

func NewUserService(userRepository user.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		userRepository: userRepository,
	}
}

func (u *UserServiceImpl) FindByID(id uint64) (*model.User, error) {
	return u.userRepository.FindByID(id)
}

func (u *UserServiceImpl) IzExist(username string) (bool, error) {
	exist, err := u.userRepository.IzExist(username)
	return exist, err
}

func (u *UserServiceImpl) Create(user *model.User) error {
	return u.userRepository.Create(user)
}

func (u *UserServiceImpl) Login(username string, password string) (*model.User, error) {
	return u.UserServiceImpl.Login(username, password)
}
