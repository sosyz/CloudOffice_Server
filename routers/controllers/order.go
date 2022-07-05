package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sonui.cn/cloudprint/models"
	"sonui.cn/cloudprint/services"
	"sonui.cn/cloudprint/utils"
	payjs2 "sonui.cn/cloudprint/utils/payjs"
	"strconv"
	"time"
)

func OrderList(c *gin.Context) {
	// 取cookie
	userID, err := c.Request.Cookie("openid")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    3,
			"message": "openid is required",
		})
		return
	}
	if list, err := services.GetOrderOverviewList(userID.Value); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    5,
			"message": "get order list error",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"list": list,
		})
	}

}

func OrderDetail(c *gin.Context) {
	var orderID int64
	// 字符串转为int64
	if orderID, err := strconv.ParseInt(c.PostForm("orderID"), 10, 64); err != nil || orderID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    3,
			"message": "orderID is required",
		})
	}
	order, err := services.OrderInfo(orderID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    err.GetErrorCode(),
			"message": "get order info error" + err.GetErrorMsg(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":  0,
			"order": order,
		})
	}
}

// OrderMerge 生成一个包含文件ID数组的订单
func OrderMerge(c *gin.Context) {
	req := c.PostForm("option")
	if req == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    3,
			"message": "option is required",
		})
		return
	}

	// 解析文件列表
	var fids []string
	err := json.Unmarshal([]byte(req), &fids)
	if err != nil || len(fids) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    5,
			"message": "option is error",
		})
		return
	}
	fmt.Printf("fids: %v\n", fids)
	ans := 0

	// 检查文件是否存在 存在计算页数和
	for _, fid := range fids {
		tFid, _ := strconv.ParseInt(fid, 10, 64)
		tFile := models.File{Fid: tFid}
		if ok, err := tFile.Exist(); !ok || err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    5,
				"message": fmt.Sprintf("fid[%v] is not exist", fid),
			})
			return // 只要有一个文件不存在，就返回
		} else {
			// TODO: 检查文件是否已使用
			var info models.FileInfo
			err = json.Unmarshal([]byte(tFile.Info), &info)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    5,
					"message": fmt.Sprintf("fileID[%v] info is error", fid),
				})
				return
			}
			fmt.Printf("fid: %v, PageNum: %+v, info:%v\n", tFile.Fid, info.PageNum, tFile.Info)
			ans += info.PageNum
		}
	}

	// 生成订单
	var fl []int64
	for _, fid := range fids {
		tfl, _ := strconv.ParseInt(fid, 10, 64)
		fl = append(fl, tfl)
	}
	order := models.Order{
		ID:        utils.OrderSF.GetId(),
		FileList:  fl,
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
			"orderID": strconv.FormatInt(order.ID, 10),
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

	data := payjs2.Create(order.TotalFee, order.ID)

	// 商户号
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

// OrderPayNotify 订单支付回调
func OrderPayNotify(c *gin.Context) {
	var data payjs2.NotifyDataObj
	if err := c.ShouldBindJSON(&data); err != nil {
		c.String(http.StatusBadRequest, "failure")
		return
	}

	if ok, err := payjs2.CheckSign(data, utils.Config.Pay.Key); err != nil {
		return
	} else if ok {
		// 更新订单状态
		order := models.Order{
			ID:     data.OutTradeNo,
			Status: models.OrderStatusPayed,
		}

		if err := order.Save(); err != nil {
			c.String(http.StatusBadRequest, "failure")
		} else {
			c.String(http.StatusOK, "success")
		}
	} else {
		c.String(http.StatusBadRequest, "failure")
	}
}

func OrderFileRepeatRead(c *gin.Context) {

}
