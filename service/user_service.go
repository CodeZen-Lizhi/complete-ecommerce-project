package service

import (
	"ecommerce/model"
	"ecommerce/repository"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
	userRepository repository.UserRepository
}

// 确保实现类实现了接口
var _ UserService = (*UserServiceImpl)(nil)

func NewUserService(userRepository repository.UserRepository) UserService {
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
	if user.Password != "" {
		hashed, err := hashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashed
	}
	return u.userRepository.Create(user)
}

func (u *UserServiceImpl) Login(username string, password string) (*model.User, error) {
	user, err := u.userRepository.FindByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}
	if isBcryptHash(user.Password) {
		if err := comparePassword(user.Password, password); err != nil {
			return nil, errors.New("用户名或密码错误")
		}
		return user, nil
	}
	// 兼容历史明文密码，并尝试升级为哈希
	if user.Password != password {
		return nil, errors.New("用户名或密码错误")
	}
	if hashed, hashErr := hashPassword(password); hashErr == nil {
		user.Password = hashed
		_ = u.userRepository.Update(user)
	}
	return user, nil
}

func hashPassword(raw string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func comparePassword(hashed, raw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(raw))
}

func isBcryptHash(password string) bool {
	return strings.HasPrefix(password, "$2a$") || strings.HasPrefix(password, "$2b$") || strings.HasPrefix(password, "$2y$")
}
