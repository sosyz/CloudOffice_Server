package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sonui.cn/cloudprint/models"
	"sonui.cn/cloudprint/services"
	"sonui.cn/cloudprint/utils"
	"sonui.cn/cloudprint/utils/log"
	payjs2 "sonui.cn/cloudprint/utils/payjs"
	"strconv"
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
	user, _ := models.GetUserFromRequest(c)

	req := c.PostForm("option")
	if req == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    3,
			"message": "option is required",
		})
		return
	}
	// 解析文件列表
	var fids []int64
	if err := json.Unmarshal([]byte(req), &fids); err != nil || len(fids) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    5,
			"message": "option is error",
		})
		return
	}
	// 返回结果
	order, err := services.OrderCreate(user.Id, fids)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    201,
			"message": "create order error",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"orderID": strconv.FormatInt(order.ID, 10), // 防止前端处理精度丢失
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
	// 输出请求包体
	body := map[string]string{}
	// 解析www-form请求表单
	if err := c.Bind(&body); err != nil {
		log.Debug("OrderPayNotify", fmt.Sprint(err))
		c.String(http.StatusBadRequest, "failure")
	}

	if payjs2.Sign(body, utils.Config.Pay.Key) == body["sign"] {
		log.Debug("OrderPayNotify", "verify sign success")
		// 更新订单状态
		order := models.Order{
			ID:     0,
			Status: models.OrderStatusPayed,
		}

		if err := order.Save(); err != nil {
			c.String(http.StatusBadRequest, "failure")
		} else {
			c.String(http.StatusOK, "success")
		}
	} else {
		log.Debug("OrderPayNotify", fmt.Sprintf("%v", body))
		c.String(http.StatusBadRequest, "failure")
		return
	}
}

func OrderFileRepeatRead(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
