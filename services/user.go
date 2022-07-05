package services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sonui.cn/cloudprint/models"
	"sonui.cn/cloudprint/utils"
	"sonui.cn/cloudprint/utils/log"
	"sonui.cn/cloudprint/vo"
	"time"
)

func GetUserInfo(openid string) (*models.User, *log.ErrorInfo) {
	user := models.User{
		OpenId: openid,
	}
	err := user.Find()
	if err != nil {
		return nil, log.NewError(202, err.Error())
	}
	return &user, nil
}

func Login(code, from string) (*vo.TokenVo, *log.ErrorInfo) {
	var auth = &vo.TokenVo{}
	switch from {
	case "wechat":
		// 微信登录
		// 获取配置信息
		WXAppID := utils.Config.Wechat.Appid
		WXAppSecret := utils.Config.Wechat.Secret

		if WXAppID == "" || WXAppSecret == "" {
			return nil, log.NewError(1, "WXAppID or WXAppSecret is empty")
		}

		//发送http请求
		resp, err := http.Get("https://api.weixin.qq.com/sns/jscode2session?appid=" + WXAppID + "&secret=" + WXAppSecret + "&js_code=" + code + "&grant_type=authorization_code")
		if err != nil {
			return nil, log.NewError(102, err.Error())
		}

		defer utils.HttpClose(resp.Body)

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, log.NewError(102, err.Error())
		}

		// 对resp转为JSON
		var p vo.Code2SessionVo
		err = json.Unmarshal(b, &p)
		if err != nil {
			return nil, log.NewError(101, err.Error())
		}

		if p.ErrCode != 0 {
			var msg string
			switch p.ErrCode {
			case -1:
				msg = "wx system busier"
			case 40029:
				msg = "invalid code"
			case 45011:
				msg = "too many request"
			case 40226:
				msg = "high-risk user"
			}
			return nil, log.NewError(101, msg)
		}

		auth.Openid = p.Openid
		auth.Session = p.SessionKey
	}

	cache := models.Cache{
		Key:      auth.Session,
		Value:    auth,
		ExpireAt: time.Now().Add(time.Hour * 24 * 365),
	}

	if err := cache.Create(); err != nil {
		return nil, log.NewError(201, "save session error")
	}
	return auth, nil
}

func SetUserInfo(openid, name, phone, address string) *log.ErrorInfo {
	user := models.User{
		OpenId:  openid,
		Name:    name,
		Phone:   phone,
		Address: address,
	}

	err := user.Save()
	if err != nil {
		return log.NewError(201, err.Error())
	}
	return nil
}
