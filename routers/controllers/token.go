package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sonui.cn/cloudprint/models"
)

type Token struct {
	Openid  string `form:"openid" json:"openid" uri:"openid" xml:"openid" binding:"required"`
	Session string `form:"session" json:"session" uri:"session" xml:"session" binding:"required"`
}

func TokenHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 特判路径是否为Login接口 是就放过
		// 不可以模糊匹配
		if c.Request.URL.Path == "/api/v1/user/login" {
			c.Next()
			return
		}

		// 获取header里的token值，没有话，就通过 c.Abort() 方法取消请求的继续进行，从而抛出异常
		var token Token
		if err := c.Bind(&token); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code": 1,
				"msg":  "get session error",
			})
			return
		}
		// 判断token是否合法
		if token.Openid == "" || token.Session == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code": 1,
				"msg":  "session is invalid",
			})
			return
		}

		// 判断token是否过期
		// 如果不一致，就通过 c.Abort() 方法取消请求的继续进行，从而抛出异常
		var value = &Token{}
		cache := models.Cache{
			Key:   token.Session,
			Value: value,
		}

		err := cache.GetValue()
		// 空记录
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code": 3,
				"msg":  "session is invalid",
			})
			return
		}

		// 记录不一致
		if !(value.Openid == token.Openid) {
			// 无法验证token
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code": 2,
				"msg":  "cant verify session",
			})
			return
		}

		// 让程序继续正常运行
		c.Next()
	}
}
