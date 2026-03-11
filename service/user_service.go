package service

import (
	"context"
	"ecommerce/internal/mysql"
	"ecommerce/model"
	"ecommerce/repository"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService 用户服务接口
type UserService interface {
	FindByID(id uint64) (*model.User, error)
	IsExist(username string) (bool, error)
	Create(user *model.User) error
	Register(ctx context.Context, user *model.User) error
	Login(username string, password string) (*model.User, error)
}

// UserServiceImpl 用户服务实现
type UserServiceImpl struct {
	userRepository repository.UserRepository
}

// 编译时检查接口实现
var _ UserService = (*UserServiceImpl)(nil)

// ErrUserAlreadyExists 表示用户名已存在。
var ErrUserAlreadyExists = errors.New("user already exists")

// NewUserService 创建用户服务实例
func NewUserService(userRepository repository.UserRepository) UserService {
	return &UserServiceImpl{
		userRepository: userRepository,
	}
}

// FindByID 根据 ID 查询用户
func (u *UserServiceImpl) FindByID(id uint64) (*model.User, error) {
	return u.userRepository.FindByID(id)
}

// IsExist 判断用户名是否存在
func (u *UserServiceImpl) IsExist(username string) (bool, error) {
	exist, err := u.userRepository.IsExist(username)
	return exist, err
}

// Create 创建用户并处理密码加密
func (u *UserServiceImpl) Create(user *model.User) error {
	if user == nil {
		return errors.New("用户不能为空")
	}
	if user.Password != "" {
		hashed, err := hashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashed
	}
	return u.userRepository.Create(user)
}

// Register 在同一事务内执行用户存在性校验与创建
func (u *UserServiceImpl) Register(ctx context.Context, user *model.User) error {
	if user == nil {
		return errors.New("用户不能为空")
	}
	return mysql.Transaction(ctx, func(tx *gorm.DB) error {
		repo := u.userRepository.WithDB(tx)
		exist, err := repo.IsExist(user.Username)
		if err != nil {
			return err
		}
		if exist {
			return ErrUserAlreadyExists
		}
		if user.Password != "" {
			hashed, err := hashPassword(user.Password)
			if err != nil {
				return err
			}
			user.Password = hashed
		}
		return repo.Create(user)
	})
}

// Login 校验用户名密码并返回用户信息
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

// hashPassword 使用 bcrypt 对明文密码进行哈希
func hashPassword(raw string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// comparePassword 比较哈希密码和明文密码是否匹配
func comparePassword(hashed, raw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(raw))
}

// isBcryptHash 判断密码字段是否为 bcrypt 格式
func isBcryptHash(password string) bool {
	return strings.HasPrefix(password, "$2a$") || strings.HasPrefix(password, "$2b$") || strings.HasPrefix(password, "$2y$")
}
