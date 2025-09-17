package model

import (
	"time"
)

// Dept 部门表结构体
type Dept struct {
	Id         int64      `json:"id" gorm:"column:id;primary_key;autoIncrement" comment:"主键ID"`
	ParentId   *int64     `json:"parent_id" gorm:"column:parent_id" comment:"父部门ID"`
	Name       *string    `json:"name" gorm:"column:name" comment:"部门名称"`
	Order      *int       `json:"order" gorm:"column:order" comment:"排序序号"`
	Level      *int       `json:"level" gorm:"column:level" comment:"部门级别"`
	Status     *string    `json:"status" gorm:"column:status" comment:"状态（例如：0-禁用，1-启用）"`
	DelFlag    *string    `json:"del_flag" gorm:"column:del_flag" comment:"删除标志（例如：0-未删除，1-已删除）"`
	CreateId   *int64     `json:"create_id" gorm:"column:create_id" comment:"创建人ID"`
	CreateTime *time.Time `json:"create_time" gorm:"column:create_time;type:datetime" comment:"创建时间"`
	UpdateId   *int64     `json:"update_id" gorm:"column:update_id" comment:"更新人ID"`
	UpdateTime *time.Time `json:"update_time" gorm:"column:update_time;type:datetime" comment:"更新时间"`
}

// TableName 指定表名
func (Dept) TableName() string {
	return "s_dept"
}
