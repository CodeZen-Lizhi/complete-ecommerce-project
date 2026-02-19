package model

import (
	"ecommerce/util"
	"time"

	"gorm.io/gorm"
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
	CreateTime *time.Time `gorm:"column:create_time;type:datetime" comment:"创建时间"`
	UpdateId   *int64     `gorm:"column:update_id" comment:"更新人ID"`
	UpdateTime *time.Time `gorm:"column:update_time;type:datetime" comment:"更新时间"`
}

//字段有指针 可以存 null 不加指针就会有默认值

// TableName 自定义表名
func (u *User) TableName() string {
	return "users"
}

// BeforeCreate 创建前钩子：自动填充创建时间和更新时间
func (u *User) BeforeCreate(tx *gorm.DB) error {
	now := time.Now() // 获取当前时间（time.Time类型）
	delFlag := "0"
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

// UserVo 定义用户请求/响应传输对象。
type UserVo struct {
	//`,string` 表示在 JSON 序列化和反序列化过程中，该 `int64` 类型的字段将被当作字符串处理。
	//这通常用于防止 JavaScript 等语言在处理大整数时出现精度丢失的问题
	ID       int64  `json:"id,string"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Age      int    `json:"age"`
}
