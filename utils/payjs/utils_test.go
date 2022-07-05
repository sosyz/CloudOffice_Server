package payjs

import (
	"testing"
)

// 测试签名验证
func TestCheckSign(t *testing.T) {
	resp := map[string]string{}
	resp["return_code"] = "1"
	resp["return_msg"] = "success"
	resp["payjs_order_id"] = "1234567890"
	resp["sign"] = Sign(resp, "1234567890")
	obj := CloseObj{
		ReturnCode:   1,
		ReturnMsg:    resp["return_msg"],
		PayJsOrderId: resp["payjs_order_id"],
		Sign:         resp["sign"],
	}
	if ok, err := CheckSign(&obj, "1234567890"); err != nil {
		t.Error(err)
	} else if !ok {
		t.Error("sign check 1 failed")
	}

	obj.Sign = "1234567890"
	if ok, err := CheckSign(&obj, "1234567890"); err != nil {
		t.Error(err)
	} else if ok {
		t.Error("sign check 2 failed")
	}
}
