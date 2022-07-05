package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sonui.cn/cloudprint/services"
	"sonui.cn/cloudprint/vo"
)

// UserGetInfo 账号信息
func UserGetInfo(client *gin.Context) {
	ck, _ := client.Request.Cookie("openid")
	res, err := services.GetUserInfo(ck.Value)
	if err != nil {
		client.JSON(
			http.StatusOK,
			gin.H{
				"code":    err.GetErrorCode(),
				"message": err.GetErrorMsg(),
			})
	} else {
		body := vo.UserInfo{
			Name:    res.Name,
			Phone:   res.Phone,
			Group:   res.Group,
			Address: res.Address,
		}
		client.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": body,
		})
	}
}

// LoginDefault 登录
func LoginDefault(client *gin.Context) {
	var login vo.LoginVo

	// 从POST参数中获得code
	if err := client.Bind(&login); err != nil {
		client.JSON(http.StatusBadRequest, gin.H{
			"code":    4,
			"message": "get code error",
		})
		return
	}

	if login.Code == "" {
		client.JSON(http.StatusBadRequest, gin.H{
			"code":    3,
			"message": "code is empty",
		})
		return
	}

	// 获取session_key
	token, err := services.Login(login.Code, login.From)
	if err != nil {
		client.JSON(
			http.StatusOK,
			gin.H{
				"code":    err.GetErrorCode(),
				"message": err.GetErrorMsg(),
			})
	} else {
		client.JSON(http.StatusOK, gin.H{
			"code":    0,
			"openid":  token.Openid,
			"session": token.Session,
		})
	}
}

func LoginQR(c *gin.Context) {

}

func LoginCheck(c *gin.Context) {

}

func LoginStatus(c *gin.Context) {

}

// UserSetInfo 设置用户信息
func UserSetInfo(c *gin.Context) {
	userInfo := vo.UserInfo{
		Name:    c.PostForm("name"),
		Phone:   c.PostForm("phone"),
		Address: c.PostForm("address"),
	}

	if userInfo.Name == "" && userInfo.Phone == "" && userInfo.Address == "" {
		c.JSON(
			http.StatusOK, gin.H{
				"code":    302,
				"message": "name, phone, address is empty",
			})
		return
	}
	ck, _ := c.Request.Cookie("openid")
	err := services.SetUserInfo(ck.Value, userInfo.Name, userInfo.Phone, userInfo.Address)

	c.JSON(
		http.StatusOK, gin.H{
			"code":    err.GetErrorCode(),
			"message": err.GetErrorMsg(),
		})
}
