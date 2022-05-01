package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"sonui.cn/cloudprint/pkg/utils"
	"strconv"
	"time"
)

func OrderList(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func OrderDetail(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func OrderSynthesis(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func OrderCancel(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// OrderPayInfo 订单支付信息
func OrderPayInfo(c *gin.Context) {
	// 商户号
	mchid := utils.GetEnvDefault("mchid", "")

	h := md5.New()
	// 取时间戳
	t := time.Now().Unix()
	h.Write([]byte(strconv.FormatInt(t, 10) + "&q^x*9@mLN3#brTTJ"))
	nonceStr := hex.EncodeToString(h.Sum(nil))

	data := map[string]string{
		"mchid":        mchid,
		"total_fee":    "1",
		"out_trade_no": "123123123123",
		"nonceStr":     nonceStr,
	}

	// PAYJS通信密钥
	key := utils.GetEnvDefault("payjs_key", "cnRAeFhI6f6f2JRd")
	sign := utils.PayJsSign(data, key)
	data["sign"] = sign

	c.JSON(200, gin.H{
		"code": 0,
		"opt":  data,
	})
}

func OrderPayStatus(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
