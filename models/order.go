package models

import (
	"encoding/json"
	"time"
)

type Order struct {
	ID         int64   `gorm:"primary_key"`
	FileList   []int64 `gorm:"-"`
	Files      string  `gorm:"column:file_list"`
	Status     int
	UserID     string
	TotalFee   int
	CreatedAt  time.Time
	PayAt      time.Time `gorm:"default:null"`
	NotifyInfo NotifyInfo
}

// NotifyInfo 回调信息
type NotifyInfo struct {
	ReturnCode    int    //	Y	1：支付成功
	TotalFee      int    //	Y	金额。单位：分
	OutTradeNo    string //	Y	用户端自主生成的订单号
	PayJsOrderId  string //	Y	PAYJS 订单号
	TransactionId string //	Y	微信用户手机显示订单号
	TimeEnd       string //	Y	支付成功时间
	Openid        string //	Y	用户OPENID标示，本参数没有实际意义，旨在方便用户端区分不同用户
	Attach        string //	N	用户自定义数据
	MchId         string //	Y	商户号
	PayType       string //	N	支付类型。微信订单不返回该字段；支付宝订单返回：alipay
	Sign          string //	Y	数据签名 详见签名算法
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
	return DB.Update(o).Error
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
	return DB.Where("id = ?", o.ID).First(o).Error
}
