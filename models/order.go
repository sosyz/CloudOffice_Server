package models

import (
	"encoding/json"
	"sonui.cn/cloudprint/utils/payjs"
	"time"
)

type Order struct {
	ID         int64   `gorm:"primary_key"`
	FileList   []int64 `gorm:"-"`
	Files      string  `gorm:"column:file_list" json:"-"`
	Status     int
	UserID     string
	TotalFee   int
	CreatedAt  time.Time
	PayAt      time.Time           `gorm:"default:null"`
	NotifyInfo payjs.NotifyDataObj `json:"-"`
}

const (
	// OrderStatusWaitPay 订单待支付
	OrderStatusWaitPay = iota
	// OrderStatusPayed 支付成功
	OrderStatusPayed
	// OrderStatusCancel 取消订单
	OrderStatusCancel
	// OrderStatusRefund 已退款
	OrderStatusRefund
	// OrderStatusTimeout 超时
	OrderStatusTimeout
	// OrderStatusError 错误
	OrderStatusError
)

func (o *Order) Save() error {
	res, _ := json.Marshal(o.FileList)
	o.Files = string(res)
	return DB.Where("id = ?", o.ID).Save(o).Error
}

func (o *Order) Create() error {
	res, _ := json.Marshal(o.FileList)
	o.Files = string(res)
	return DB.Create(o).Error
}

func (o *Order) Delete() error {
	return DB.Delete(o).Error
}

func (o *Order) Find() error {
	if err := DB.Where("id = ?", o.ID).First(o).Error; err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(o.Files), &o.FileList); err != nil {
		return err
	}

	return nil
}

// OrderList 根据用户ID查找订单列表
func OrderList(userID string) ([]*Order, error) {
	var orders []*Order
	err := DB.Where("user_id = ?", userID).Find(&orders).Error
	for _, v := range orders {
		if err := json.Unmarshal([]byte(v.Files), &v.FileList); err != nil {
			return nil, err
		}
	}
	return orders, err
}
