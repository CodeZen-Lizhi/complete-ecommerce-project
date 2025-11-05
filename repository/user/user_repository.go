package user

import (
	"ecommerce/internal/logger"
	"ecommerce/internal/mysql"
	"ecommerce/model"
	"gorm.io/gorm"
)

// UserRepository 定义接口
type UserRepository interface {
	FindByID(id uint64) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
	IzExist(username string) (bool, error)
	Login(username string, password string) (*model.User, error)
}

// 实现接口
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 构造函数注入
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
	tx := r.db.Create(user)
	return tx.Error
}

func (r *userRepository) Update(user *model.User) error {
	tx := r.db.Updates(user)
	return tx.Error
}

func (r *userRepository) FindByID(id uint64) (*model.User, error) {
	var user model.User
	tx := r.db.First(&user, id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	log := logger.GetLogger()
	log.Debug("查询到用户", "user", user)
	return &user, tx.Error
}

func (r *userRepository) IzExist(username string) (bool, error) {
	var user model.User
	tx := r.db.Where("username = ? and del_flag = ?", username, "0").First(&user)
	if tx.Error != nil {
		return false, tx.Error
	}
	return tx.RowsAffected > 0, tx.Error
}

func (r *userRepository) Login(username string, password string) (*model.User, error) {
	var user model.User
	tx := r.db.Where("username = ? and password = ? and del_flag = ?", username, password, "0").First(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &user, tx.Error
}
