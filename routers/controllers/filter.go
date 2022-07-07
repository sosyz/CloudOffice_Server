package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sonui.cn/cloudprint/models"
)

var whiteList = []string{
	"/api/v1/user/login",
	"/api/v1/user/register",
	"/api/v1/order/pay/notify",
	"/api/v1/ws",
}

func TokenHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 特判路径是否为白名单接口 是就放过
		// 不可以模糊匹配
		for _, path := range whiteList {
			if c.Request.URL.Path == path {
				c.Next()
				return
			}
		}
		_, err := models.GetUserFromRequest(c)

		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    err.GetErrorCode(),
				"message": err.GetErrorMsg(),
			})
			c.Abort()
			return
		}

		// 让程序继续正常运行
		c.Next()
	}
}
