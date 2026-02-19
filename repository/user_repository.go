package repository

import (
	"ecommerce/internal/logger"
	"ecommerce/internal/mysql"
	"ecommerce/model"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"

	"gorm.io/gorm"
)

// UserRepository 定义接口
type UserRepository interface {
	FindByID(id uint64) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
	IzExist(username string) (bool, error)
}

// UserRepositoryImpl 实现接口
type UserRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository 构造函数注入
func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{
		db: mysql.DB,
	}
}

// 此段代码可以确保结构体实现了接口的所有方法，否则编译会出错
var _ UserRepository = (*UserRepositoryImpl)(nil)

// FindByEmail 根据邮箱查询用户。
func (r *UserRepositoryImpl) FindByEmail(email string) (*model.User, error) {
	var user model.User
	tx := r.db.Where("email = ? and del_flag = ?", email, "0").First(&user)
	if tx.Error != nil {
		return nil, wrapRepoError("UserRepositoryImpl.FindByEmail", tx.Error)
	}
	return &user, nil
}

// Create 创建用户记录。
func (r *UserRepositoryImpl) Create(user *model.User) error {
	tx := r.db.Create(user)
	return wrapRepoError("UserRepositoryImpl.Create", tx.Error)
}

// Update 更新用户记录。
func (r *UserRepositoryImpl) Update(user *model.User) error {
	tx := r.db.Updates(user)
	return wrapRepoError("UserRepositoryImpl.Update", tx.Error)
}

// FindByID 根据ID查询用户。
func (r *UserRepositoryImpl) FindByID(id uint64) (*model.User, error) {
	var user model.User
	tx := r.db.First(&user, id)
	if tx.Error != nil {
		return nil, wrapRepoError("UserRepositoryImpl.FindByID", tx.Error)
	}
	log := logger.GetLogger()
	log.Debug("查询到用户", "user", user)
	return &user, nil
}

// IzExist 判断用户名是否已存在。
func (r *UserRepositoryImpl) IzExist(username string) (bool, error) {
	var user model.User
	tx := r.db.Where("username = ? and del_flag = ?", username, "0").First(&user)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, wrapRepoError("UserRepositoryImpl.IzExist", tx.Error)
	}
	return tx.RowsAffected > 0, nil
}

// FindByUsername 根据用户名查询用户。
func (r *UserRepositoryImpl) FindByUsername(username string) (*model.User, error) {
	var user model.User
	tx := r.db.Where("username = ? and del_flag = ?", username, "0").First(&user)
	if tx.Error != nil {
		return nil, wrapRepoError("UserRepositoryImpl.FindByUsername", tx.Error)
	}
	return &user, nil
}

// wrapRepoError 包装仓储层错误并附带代码行信息。
func wrapRepoError(operation string, err error) error {
	if err == nil {
		return nil
	}
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return fmt.Errorf("%s failed: %w", operation, err)
	}
	return fmt.Errorf("%s failed (%s:%d): %w", operation, filepath.Base(file), line, err)
}
