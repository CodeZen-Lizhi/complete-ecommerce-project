package model

import (
	"time"
)

// Dict 字典主表结构体
type Dict struct {
	Id          int64      `json:"id" gorm:"column:id;primary_key;autoIncrement" comment:"主键ID"`
	DictName    string     `json:"dict_name" gorm:"column:dict_name;not null;size:128" comment:"字典名称"`
	DictCode    string     `json:"dict_code" gorm:"column:dict_code;not null;size:64;uniqueIndex:uk_dict_code" comment:"字典编码（唯一标识）"`
	Description *string    `json:"description" gorm:"column:description;size:512" comment:"字典描述"`
	Type        *string    `json:"type" gorm:"column:type;type:char" comment:"字典类型（例如：0-系统字典，1-业务字典）"`
	Status      *string    `json:"status" gorm:"column:status;type:char;default:'1'" comment:"状态（0-禁用，1-启用）"`
	DelFlag     *string    `json:"del_flag" gorm:"column:del_flag;type:char;default:'0'" comment:"删除标志（0-未删除，1-已删除）"`
	CreateId    *int64     `json:"create_id" gorm:"column:create_id" comment:"创建人ID"`
	CreateTime  *time.Time `json:"create_time" gorm:"column:create_time;default:CURRENT_TIMESTAMP" comment:"创建时间"`
	UpdateId    *int64     `json:"update_id" gorm:"column:update_id" comment:"更新人ID"`
	UpdateTime  *time.Time `json:"update_time" gorm:"column:update_time;default:CURRENT_TIMESTAMP;autoUpdateTime" comment:"更新时间"`
}

// TableName 指定表名
func (Dict) TableName() string {
	return "s_dict"
}
