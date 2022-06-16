package services

import (
	"encoding/json"
	"fmt"
	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
	"io/ioutil"
	"net/http"
	"sonui.cn/cloudprint/models"
	"sonui.cn/cloudprint/utils"
	"sonui.cn/cloudprint/vo"
	"time"
)

func GetUserInfo(openid string) (int, string, *models.User) {
	user := models.User{
		OpenId: openid,
	}
	err := user.Find()
	if err != nil {
		return 202, err.Error(), nil
	}
	return 0, "", &user
}

// CreatTmpKey 创建临时密钥
func CreatTmpKey(openid string) (int, string, any) {
	if utils.Config.QCloud.SecretId == "" || utils.Config.QCloud.SecretKey == "" {
		return 1, "SecretId or SecretKey is empty", nil
	}

	c := sts.NewClient(
		utils.Config.QCloud.SecretId,
		utils.Config.QCloud.SecretKey,
		nil,
	)

	// 获取存储桶信息
	if utils.Config.QCloud.Region == "" || utils.Config.QCloud.Appid == "" || utils.Config.QCloud.Bucket == "" {
		return 1, "COS option is empty", nil
	}

	// 查询缓存
	var value = &sts.CredentialResult{}
	cache := models.Cache{
		Key:   openid + "_tmpKey",
		Value: value,
	}

	if err := cache.Find(); err == nil {
		return 0, "", cache.Value
	}

	opt := &sts.CredentialOptions{
		DurationSeconds: int64(time.Hour.Seconds()),
		Region:          utils.Config.QCloud.Region,
		Policy: &sts.CredentialPolicy{
			Statement: []sts.CredentialPolicyStatement{
				{
					// 密钥的权限列表。简单上传和分片需要以下的权限，其他权限列表请看 https://cloud.tencent.com/document/product/436/31923
					Action: []string{
						// 简单上传
						"name/cos:PostObject",
						"name/cos:PutObject",
						// 分片上传
						"name/cos:InitiateMultipartUpload",
						"name/cos:ListMultipartUploads",
						"name/cos:ListParts",
						"name/cos:UploadPart",
						"name/cos:CompleteMultipartUpload",
					},
					Effect: "allow",
					Resource: []string{
						"qcs::cos:" + utils.Config.QCloud.Region + ":uid/" + utils.Config.QCloud.Appid + ":" + utils.Config.QCloud.Bucket + "/" + openid + "/*",
					},
				},
			},
		},
	}

	// 请求临时密钥
	res, err := c.GetCredential(opt)
	if err != nil {
		return 105, err.Error(), nil
	} else {
		// 缓存密钥
		// 时间戳转时间
		timeAt := time.Unix(int64(res.ExpiredTime), 0)
		cache := models.Cache{
			Key:      openid + "_tmpKey",
			Value:    res,
			ExpireAt: timeAt,
		}
		_ = cache.Create()

		return 0, "", res
	}
}

func Login(code string) (int, string, *vo.TokenVo) {
	// 获取配置信息
	WXAppID := utils.Config.Wechat.Appid
	WXAppSecret := utils.Config.Wechat.Secret

	if WXAppID == "" || WXAppSecret == "" {
		return 1, "WXAppID or WXAppSecret is empty", nil
	}

	//发送http请求
	resp, err := http.Get("https://api.weixin.qq.com/sns/jscode2session?appid=" + WXAppID + "&secret=" + WXAppSecret + "&js_code=" + code + "&grant_type=authorization_code")
	if err != nil {
		return 101, err.Error(), nil
	}

	defer utils.HttpClose(resp.Body)

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 106, err.Error(), nil
	}
	fmt.Printf("%s\n", b)
	// 对resp转为JSON
	var p vo.Code2SessionVo
	err = json.Unmarshal(b, &p)
	if err != nil {
		return 102, err.Error(), nil
	}
	if p.ErrCode != 0 {
		switch p.ErrCode {
		case -1:
			return 101, "wx system busier", nil
		case 40029:
			return 101, "invalid code", nil
		case 45011:
			return 101, "too many request", nil
		case 40226:
			return 101, "high-risk user", nil
		}
	}
	//设置登录信息缓存
	token := &vo.TokenVo{
		Session: p.SessionKey,
		Openid:  p.Openid,
	}

	fmt.Printf("%+v\n", token)
	cache := models.Cache{
		Key:      p.SessionKey,
		Value:    token,
		ExpireAt: time.Now().Add(time.Hour * 24 * 365),
	}

	if err = cache.Create(); err != nil {
		return 201, "cache save error", nil
	}
	return 0, "", token
}

func SetUserInfo(openid, name, phone, address string) (int, string) {
	user := models.User{
		OpenId:  openid,
		Name:    name,
		Phone:   phone,
		Address: address,
	}

	err := user.Save()
	if err != nil {
		return 201, err.Error()
	}
	return 0, "success"
}
