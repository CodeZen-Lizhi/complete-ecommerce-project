package model

import (
	"time"
)

// User 用户模型
type User struct {
	ID         uint       `gorm:"primarykey" json:"id"`
	Username   string     `gorm:"size:50;uniqueIndex" json:"username"`
	Email      string     `gorm:"size:100;uniqueIndex" json:"email"`
	Password   string     `gorm:"size:100" json:"-"` // 不返回密码
	Age        int        `gorm:"default:0" json:"age"`
	DelFlag    *string    `json:"del_flag" gorm:"column:del_flag" comment:"删除标志（例如：0-未删除，1-已删除）"`
	CreateId   *int64     `json:"create_id" gorm:"column:create_id" comment:"创建人ID"`
	CreateTime *time.Time `json:"create_time" gorm:"column:create_time" comment:"创建时间"`
	UpdateId   *int64     `json:"update_id" gorm:"column:update_id" comment:"更新人ID"`
	UpdateTime *time.Time `json:"update_time" gorm:"column:update_time" comment:"更新时间"`
}

// TableName 自定义表名
func (User) TableName() string {
	return "users"
}
