package controllers

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sonui.cn/cloudprint/services"
	"sonui.cn/cloudprint/vo"
)

// GetUserInfo 账号信息
func GetUserInfo(client *gin.Context) {
	ck, _ := client.Request.Cookie("openid")
	code, err, res := services.GetUserInfo(ck.Value)
	if err != "" {
		client.JSON(
			http.StatusOK,
			gin.H{
				"code":    code,
				"message": err,
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

// CreatTmpKey 获取临时密钥
func CreatTmpKey(client *gin.Context) {
	// 前头验证过了这里忽略可能出现的错误
	openid, _ := client.Request.Cookie("openid")

	//openid取其md5
	data := []byte(openid.Value)
	has := md5.Sum(data)
	//将[]byte转成16进制
	openid.Value = fmt.Sprintf("%x", has)
	code, err, res := services.CreatTmpKey(openid.Value)

	if err != "" {
		client.JSON(
			http.StatusOK,
			gin.H{
				"code":    code,
				"message": err,
			})
	} else {
		client.JSON(
			http.StatusOK,
			gin.H{
				"code": 0,
				"data": res,
			})
	}

}

// Login 登录
func Login(client *gin.Context) {
	var form vo.LoginVo

	// 从POST参数中获得code
	if err := client.Bind(&form); err != nil {
		client.JSON(http.StatusBadRequest, gin.H{
			"code":    4,
			"message": "get code error",
		})
		return
	}

	if form.Code == "" {
		client.JSON(http.StatusBadRequest, gin.H{
			"code":    3,
			"message": "code is empty",
		})
		return
	}

	// 获取session_key
	code, err, token := services.Login(form.Code)
	if err != "" {
		client.JSON(
			http.StatusOK,
			gin.H{
				"code":    code,
				"message": err,
			})
	} else {
		client.JSON(http.StatusOK, gin.H{
			"code":    0,
			"openid":  token.Openid,
			"session": token.Session,
		})
	}
}

// SetUserInfo 设置用户信息
func SetUserInfo(c *gin.Context) {
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
	code, err := services.SetUserInfo(ck.Value, userInfo.Name, userInfo.Phone, userInfo.Address)

	c.JSON(
		http.StatusOK, gin.H{
			"code":    code,
			"message": err,
		})
}
