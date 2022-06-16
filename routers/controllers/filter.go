package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sonui.cn/cloudprint/models"
	"sonui.cn/cloudprint/vo"
)

func TokenHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 特判路径是否为Login接口 是就放过
		// 不可以模糊匹配
		if c.Request.URL.Path == "/api/v1/user/login" {
			c.Next()
			return
		}

		var token vo.TokenVo
		//获取cookie
		ck, err := c.Request.Cookie("openid")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code": 301,
				"msg":  "no login",
			})
			return
		}
		token.Openid = ck.Value

		ck, err = c.Request.Cookie("session")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code": 301,
				"msg":  "no login",
			})
			return
		}
		token.Session = ck.Value

		if token.Openid == "" || token.Session == "" {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code": 301,
				"msg":  "no login",
			})
			return
		}

		// 判断token是否过期
		var value = &vo.TokenVo{}
		cache := models.Cache{
			Key:   token.Session,
			Value: value,
		}

		err = cache.Find()
		// 空记录
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code": 301,
				"msg":  "no login",
			})
			return
		}

		// 记录不一致
		if !(value.Openid == token.Openid) {
			// 无法验证token
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code": 301,
				"msg":  "no login",
			})
			return
		}

		// 让程序继续正常运行
		c.Next()
	}
}
