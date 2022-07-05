package payjs

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sonui.cn/cloudprint/utils"
	"strconv"
	"strings"
	"time"
)

// buildUrl 参数编译 将map转换为x-www-form格式
// 参数: data map[string]string
func buildUrl(data map[string]string) string {
	var ret string
	for k, v := range data {
		ret += "&" + k + "=" + url.QueryEscape(v)
	}
	// 去除开头的&
	return ret[1:]
}

// Create 创建订单
// 参数: totalFee 总金额, orderID 自定义订单ID
func Create(totalFee int, orderID int64) interface{} {
	// 商户号
	mchId := utils.Config.Pay.MchId
	h := md5.New()
	// 取时间戳
	t := time.Now().Unix()
	h.Write([]byte(strconv.FormatInt(t, 10) + "&q^x*9@mLN3#brTTJ"))
	nonceStr := hex.EncodeToString(h.Sum(nil))

	data := map[string]string{
		"mchid":        mchId,
		"total_fee":    strconv.Itoa(totalFee),
		"out_trade_no": strconv.FormatInt(orderID, 10),
		"notify_url":   utils.Config.Run.Host + ":" + utils.Config.Run.Listen + "/pay/notify",
		"nonceStr":     nonceStr,
	}

	// 签名
	data["sign"] = Sign(data, utils.Config.Pay.Key)
	return data
}

// Status 查询订单状态
// 参数: orderID PayJS订单ID
func Status(orderID int64) (bool, error) {
	data := map[string]string{
		"payjs_order_id": strconv.FormatInt(orderID, 10),
	}
	// 签名
	data["sign"] = Sign(data, utils.Config.Pay.Key)

	// 发送Post请求
	resp, _ := http.Post("https://payjs.cn/api/check", "application/x-www-form-urlencoded", strings.NewReader(buildUrl(data)))
	defer utils.HttpClose(resp.Body)
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	// 解析返回数据
	var ret NotifyDataObj
	if err := json.Unmarshal(b, &ret); err != nil {
		return false, err
	}
	return ret.Status == 1, nil
}

// Close 关闭订单
// 参数: orderID PayJS订单ID
func Close(orderID int64) error {
	data := map[string]string{
		"payjs_order_id": strconv.FormatInt(orderID, 10),
	}
	// 签名
	data["sign"] = Sign(data, utils.Config.Pay.Key)

	// 发送Post请求
	resp, _ := http.Post("https://payjs.cn/api/close", "application/x-www-form-urlencoded", strings.NewReader(buildUrl(data)))
	defer utils.HttpClose(resp.Body)
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 解析返回数据
	var ret CloseObj
	if err := json.Unmarshal(b, &ret); err != nil {
		return err
	}
	if ret.ReturnCode != 1 {
		return fmt.Errorf("%s", ret.ReturnMsg)
	}
	return nil
}

// Reverse 取消订单
// 参数: orderID PayJS订单ID
func Reverse(orderID int64) error {
	data := map[string]string{
		"payjs_order_id": strconv.FormatInt(orderID, 10),
	}
	// 签名
	data["sign"] = Sign(data, utils.Config.Pay.Key)

	// 发送Post请求
	resp, _ := http.Post("https://payjs.cn/api/reverse", "application/x-www-form-urlencoded", strings.NewReader(buildUrl(data)))
	defer utils.HttpClose(resp.Body)
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 解析返回数据
	var ret ReverseObj
	if err := json.Unmarshal(b, &ret); err != nil {
		return err
	}
	if ret.ReturnCode != 1 {
		return fmt.Errorf("%s", ret.ReturnMsg)
	}
	return nil
}

// Refund 商户退款
// 参数: orderID PayJS订单ID
func Refund(orderID int64) error {
	data := map[string]string{
		"payjs_order_id": strconv.FormatInt(orderID, 10),
	}
	// 签名
	data["sign"] = Sign(data, utils.Config.Pay.Key)

	// 发送Post请求
	resp, _ := http.Post("https://payjs.cn/api/refund", "application/x-www-form-urlencoded", strings.NewReader(buildUrl(data)))
	defer utils.HttpClose(resp.Body)
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 解析返回数据
	var ret RefundObj
	if err := json.Unmarshal(b, &ret); err != nil {
		return err
	}
	if ret.Refund == "1" {
		return fmt.Errorf("订单已退款")
	}
	if ret.ReturnCode != 1 {
		return fmt.Errorf("%s", ret.ReturnMsg)
	}
	return nil
}
