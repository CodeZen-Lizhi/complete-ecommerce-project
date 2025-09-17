package user

import (
	"ecommerce/internal/logger"
	"ecommerce/internal/mysql"
	"ecommerce/model"
	"gorm.io/gorm"
)

// UserRepository 定义接口
type UserRepository interface {
	FindByID(id uint) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
	izExist(username string) (bool, error)
}

// 实现接口
type userRepository struct {
	db *gorm.DB
}

// 构造函数注入
func NewUserRepository() UserRepository {
	return &userRepository{
		db: mysql.DB,
	}
}

// 此段代码可以确保结构体实现了接口的所有方法，否则编译会出错
var _ UserRepository = (*userRepository)(nil)

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	//mysql.DB.select("id", "email", "password").Where("email = ?", email).First(&model.User{})

	panic("implement me")
}

func (r *userRepository) Create(user *model.User) error {
	//TODO implement me
	panic("implement me")
}

func (r *userRepository) Update(user *model.User) error {
	//TODO implement me
	panic("implement me")
}

func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	tx := r.db.First(&user, id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	log := logger.GetLogger()
	log.Debug("查询到用户", "user", user)
	return &user, tx.Error
}

func (r *userRepository) izExist(username string) (bool, error) {
	var user model.User
	tx := r.db.Where("username = ?", username).First(&user)
	if tx.Error != nil {
		return false, tx.Error
	}
	return tx.RowsAffected > 0, tx.Error
}
