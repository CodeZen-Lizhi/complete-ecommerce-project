package model

import (
	"ecommerce/util"
	"gorm.io/gorm"
	"time"
)

// User 用户模型
type User struct {
	ID         int64      `gorm:"primarykey;autoIncrement:false"`
	Username   string     `gorm:"size:50;uniqueIndex"`
	Email      string     `gorm:"size:100;uniqueIndex"`
	Password   string     `gorm:"size:100"`
	Age        int        `gorm:"default:0"`
	DelFlag    *string    `gorm:"column:del_flag" comment:"删除标志（例如：0-未删除，1-已删除）"`
	CreateId   *int64     `gorm:"column:create_id" comment:"创建人ID"`
	CreateTime *time.Time ` gorm:"column:create_time;type:datetime" comment:"创建时间"`
	UpdateId   *int64     `gorm:"column:update_id" comment:"更新人ID"`
	UpdateTime *time.Time `gorm:"column:update_time;type:datetime" comment:"更新时间"`
}

//字段有指针 可以存 null 不加指针就会有默认值

// TableName 自定义表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前钩子：自动填充创建时间和更新时间
func (u *User) BeforeCreate(tx *gorm.DB) error {
	now := time.Now() // 获取当前时间（time.Time类型）
	delFlag := "1"
	u.CreateTime = &now
	u.UpdateTime = &now
	u.DelFlag = &delFlag
	u.ID = util.GenID()
	return nil
}

// BeforeUpdate 更新前钩子：自动更新更新时间
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	now := time.Now()
	u.UpdateTime = &now // 取地址赋值
	return nil
}
