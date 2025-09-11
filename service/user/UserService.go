package user

import (
	"ecommerce/model"
	"ecommerce/repository/user"
)

type UserService interface {
	FindByID(id uint) (*model.User, error)
}

type userService struct {
	userRepository user.UserRepository
}

func NewUserService(userRepository user.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (u *userService) FindByID(id uint) (*model.User, error) {
	return u.userRepository.FindByID(id)
}
