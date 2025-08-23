package model

import (
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Username  string         `gorm:"size:50;uniqueIndex" json:"username"`
	Email     string         `gorm:"size:100;uniqueIndex" json:"email"`
	Password  string         `gorm:"size:100" json:"-"` // 不返回密码
	Age       int            `gorm:"default:0" json:"age"`
	CreatedAt int64          `json:"created_at"`
	UpdatedAt int64          `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 自定义表名
func (User) TableName() string {
	return "users"
}
