package payjs

// NotifyDataObj 异步回调结构体
type NotifyDataObj struct {
	ReturnCode    int    `json:"return_code" form:"return_code"`
	Mchid         string `json:"mchid" form:"mchid"`
	OutTradeNo    int64  `json:"out_trade_no" form:"out_trade_no"`
	PayjsOrderId  string `json:"payjs_order_id" form:"payjs_order_id"`
	TransactionId string `json:"transaction_id" form:"transaction_id"`
	Status        int    `json:"status" form:"status"`
	Openid        string `json:"openid" form:"openid"`
	TotalFee      int    `json:"total_fee" form:"total_fee"`
	PaidTime      string `json:"paid_time" form:"paid_time"`
	Attach        string `json:"attach" form:"attach"`
	Sign          string `json:"sign" form:"sign"`
}

// CloseObj 关闭订单接口返回结构体
type CloseObj struct {
	ReturnCode   int    `json:"return_code" form:"return_code"`
	ReturnMsg    string `json:"return_msg" form:"return_msg"`
	PayJsOrderId string `json:"payjs_order_id" form:"payjs_order_id"`
	Sign         string `json:"sign" form:"sign"`
}

// ReverseObj 撤销订单接口返回结构体
type ReverseObj struct {
	ReturnCode   int    `json:"return_code" form:"return_code"`
	ReturnMsg    string `json:"return_msg" form:"return_msg"`
	PayJsOrderId string `json:"payjs_order_id" form:"payjs_order_id"`
	Sign         string `json:"sign" form:"sign"`
}

// RefundObj 退款接口返回结构体
type RefundObj struct {
	ReturnCode    int    `json:"return_code" form:"return_code"`
	ReturnMsg     string `json:"return_msg" form:"return_msg"`
	PayJsOrderId  string `json:"payjs_order_id" form:"payjs_order_id"`
	OutTradeNo    string `json:"out_trade_no" form:"out_trade_no"`
	TransactionId string `json:"transaction_id" form:"transaction_id"`
	Refund        string `json:"refund" form:"refund"`
	Sign          string `json:"sign" form:"sign"`
}
