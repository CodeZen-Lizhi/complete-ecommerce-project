package model

import (
	"time"
)

// DictItem 字典项表结构体
type DictItem struct {
	Id            int64      `json:"id" gorm:"column:id;primary_key;autoIncrement" comment:"主键ID"`
	DictCode      string     `json:"dict_code" gorm:"column:dict_code;not null;size:64" comment:"字典编码（关联字典主表）"`
	DictItemCode  string     `json:"dict_item_code" gorm:"column:dict_item_code;not null;size:64" comment:"字典项编码"`
	DictItemValue *string    `json:"dict_item_value" gorm:"column:dict_item_value;size:255" comment:"字典项值"`
	Status        *string    `json:"status" gorm:"column:status;type:char;default:'1'" comment:"状态（0-禁用，1-启用）"`
	Sort          *int       `json:"sort" gorm:"column:sort;default:0" comment:"排序序号"`
	DelFlag       *string    `json:"del_flag" gorm:"column:del_flag;type:char;default:'0'" comment:"删除标志（0-未删除，1-已删除）"`
	CreateId      *int64     `json:"create_id" gorm:"column:create_id" comment:"创建人ID"`
	CreateTime    *time.Time `json:"create_time" gorm:"column:create_time;default:CURRENT_TIMESTAMP" comment:"创建时间"`
	UpdateId      *int64     `json:"update_id" gorm:"column:update_id" comment:"更新人ID"`
	UpdateTime    *time.Time `json:"update_time" gorm:"column:update_time;default:CURRENT_TIMESTAMP;autoUpdateTime" comment:"更新时间"`
}

// TableName 指定表名
func (DictItem) TableName() string {
	return "s_dict_item"
}

// Indexes 定义联合唯一索引（与数据库约束uk_dict_code_item_code对应）
func (DictItem) Indexes() map[string]interface{} {
	return map[string]interface{}{
		"uk_dict_code_item_code": []string{"dict_code", "dict_item_code"},
	}
}
