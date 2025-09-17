package model

import (
	"time"
)

// Order 订单表结构体
type Order struct {
	ID         uint       `json:"id" gorm:"column:id;primaryKey;autoIncrement" comment:"订单ID"`
	UserID     int64      `json:"user_id" gorm:"column:user_id;not null" comment:"用户ID"`
	OrderNo    string     `json:"order_no" gorm:"column:order_no;size:50;uniqueIndex:idx_order_no;not null" comment:"订单号"`
	Amount     float64    `json:"amount" gorm:"column:amount;type:decimal(10,2);not null" comment:"订单总金额"`
	Status     string     `json:"status" gorm:"column:status;size:20;not null" comment:"订单状态"`
	Address    string     `json:"address" gorm:"column:address;size:255;not null" comment:"收货地址"`
	ProductID  *int64     `json:"product_id" gorm:"column:product_id" comment:"商品ID"`
	PayTime    *time.Time `json:"pay_time" gorm:"column:pay_time" comment:"支付时间"`
	DelFlag    *string    `json:"del_flag" gorm:"column:del_flag;type:char;default:'0'" comment:"删除标志（0-未删除，1-已删除）"`
	CreateId   *int64     `json:"create_id" gorm:"column:create_id" comment:"创建人ID"`
	CreateTime *time.Time `json:"create_time" gorm:"column:create_time;type:datetime" comment:"创建时间"`
	UpdateId   *int64     `json:"update_id" gorm:"column:update_id" comment:"更新人ID"`
	UpdateTime *time.Time `json:"update_time" gorm:"column:update_time;type:datetime;autoUpdateTime" comment:"更新时间"`
}

// TableName 指定表名
func (Order) TableName() string {
	return "t_order"
}

// Indexes 定义索引
func (Order) Indexes() map[string]interface{} {
	return map[string]interface{}{
		"idx_user_id":  []string{"user_id"},
		"idx_del_flag": []string{"del_flag"},
	}
}
