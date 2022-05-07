package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sonui.cn/cloudprint/models"
	"sonui.cn/cloudprint/pkg/conf"
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

// OrderMerge 生成一个包含文件ID数组的订单
func OrderMerge(c *gin.Context) {
	req := c.PostForm("option")
	if req == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    3,
			"message": "option is required",
		})
		return
	}

	// 解析文件列表
	var fids []int64
	err := json.Unmarshal([]byte(req), &fids)
	if err != nil || len(fids) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    5,
			"message": "option is error",
		})
		return
	}

	var ans int
	var info models.FileInfo
	// 检查文件是否存在 存在计算页数和
	for _, fid := range fids {
		tFile := models.File{Fid: fid}
		if !tFile.Exist() {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    5,
				"message": fmt.Sprintf("fid[%v] is not exist", fid),
			})
			return // 只要有一个文件不存在，就返回
		} else {
			// TODO: 检查文件是否已使用
			err = json.Unmarshal([]byte(tFile.Info), &info)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    5,
					"message": fmt.Sprintf("fileID[%d] info is error", fid),
				})
				return
			}
			fmt.Printf("fid: %v, PageNum: %+v\n", fid, info.PageNum)
			ans += info.PageNum
		}
	}

	// 生成订单
	order := models.Order{
		ID:        utils.OrderSF.GetId(),
		FileList:  fids,
		Status:    models.OrderStatusWaitPay,
		UserID:    c.PostForm("openid"),
		TotalFee:  ans * 30,
		CreatedAt: time.Now(),
	}

	// 写入数据库
	err = order.Create()
	// 返回结果
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    201,
			"message": "create order error",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"orderID": order.ID,
		})
	}
}

func OrderCancel(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// OrderPayInfo 订单支付信息
func OrderPayInfo(c *gin.Context) {
	orderID, err := strconv.ParseInt(c.PostForm("orderID"), 10, 64)
	if orderID == 0 || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    3,
			"message": "orderID is required",
		})
		return
	}

	// 获取订单信息
	order := models.Order{
		ID: orderID,
	}

	if err = order.Find(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    3,
			"message": "order is not exist",
		})
		return
	}
	fmt.Printf("order: %+v\n", order)

	// 商户号
	mchId := conf.Conf.Pay.MchId
	h := md5.New()
	// 取时间戳
	t := time.Now().Unix()
	h.Write([]byte(strconv.FormatInt(t, 10) + "&q^x*9@mLN3#brTTJ"))
	nonceStr := hex.EncodeToString(h.Sum(nil))

	data := map[string]string{
		"mchid":        mchId,
		"total_fee":    strconv.Itoa(order.TotalFee),
		"out_trade_no": strconv.FormatInt(orderID, 10),
		"nonceStr":     nonceStr,
	}

	// PAYJS通信密钥
	key := conf.Conf.Pay.Key
	sign := utils.PayJsSign(data, key)
	data["sign"] = sign

	c.JSON(http.StatusOK, gin.H{
		"code":   0,
		"payOpt": data,
	})
}

func OrderPayStatus(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func PayNotify(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
