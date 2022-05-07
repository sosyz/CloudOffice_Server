package controllers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"sonui.cn/cloudprint/models"
	"sonui.cn/cloudprint/pkg/conf"
	"time"

	"github.com/gin-gonic/gin"
	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
)

// TODO: 目前还是太乱细分以后再做吧

type code2Session struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

type LoginObj struct {
	Code string `form:"code" json:"code" uri:"code" xml:"code" binding:"required"`
}

// CreatTmpKey 获取临时密钥
func CreatTmpKey(client *gin.Context) {
	// 存在安全隐患 对外来参数可能有伪造 _key
	var token Token
	_ = client.Bind(&token) // 前头验证过了这里忽略

	//openid取其md5
	//data := []byte(form.Openid)
	//has := md5.Sum(data)
	//form.Openid = fmt.Sprintf("%x", has) //将[]byte转成16进制
	SecretId := conf.Conf.Secret.SecretId
	SecretKey := conf.Conf.Secret.SecretKey
	if SecretId == "" || SecretKey == "" {
		client.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "SecretId or SecretKey is empty",
		})
		return
	}

	c := sts.NewClient(
		SecretId,
		SecretKey,
		nil,
	)

	// 获取存储桶信息
	CosBucketRegion := conf.Conf.Cos.Region
	if CosBucketRegion == "" {
		client.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "CosBucket is empty",
		})
		return
	}

	CosAppid := conf.Conf.Cos.Appid
	if CosAppid == "" {
		client.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "CosAppid is empty",
		})
		return
	}

	CosBucketName := conf.Conf.Cos.Bucket
	if CosBucketName == "" {
		client.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "CosBucketName is empty",
		})
		return
	}

	// 查询缓存
	var value = &sts.CredentialResult{}
	cache := models.Cache{
		Key:   token.Openid + "_tmpKey",
		Value: value,
	}

	if err := cache.GetValue(); err == nil {
		client.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": cache.Value,
		})
		return
	}

	opt := &sts.CredentialOptions{
		DurationSeconds: int64(time.Hour.Seconds()),
		Region:          CosBucketRegion,
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
						"qcs::cos:" + CosBucketRegion + ":uid/" + CosAppid + ":" + CosBucketName + "/" + token.Openid + "/*",
					},
				},
			},
		},
	}

	// 请求临时密钥
	res, err := c.GetCredential(opt)

	if err != nil {
		client.JSON(http.StatusBadRequest, gin.H{
			"code":    105,
			"message": err,
		})
	} else {
		// 缓存密钥
		// 时间戳转时间
		timeAt := time.Unix(int64(res.ExpiredTime), 0)
		cache := models.Cache{
			Key:      token.Openid + "_tmpKey",
			Value:    res,
			ExpireAt: timeAt,
		}
		_ = cache.Create()

		client.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": res,
		})
	}
}

func Login(client *gin.Context) {
	var form LoginObj
	if err := client.Bind(&form); err != nil {
		client.JSON(http.StatusBadRequest, gin.H{
			"code":    4,
			"message": "bind option error",
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

	// 获取配置信息
	WXAppID := conf.Conf.Wechat.Appid
	WXAppSecret := conf.Conf.Wechat.Secret

	if WXAppID == "" || WXAppSecret == "" {
		client.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "WXAppID or WXAppSecret is empty",
		})
		return
	}

	//发送http请求
	resp, err := http.Get("https://api.weixin.qq.com/sns/jscode2session?appid=" + WXAppID + "&secret=" + WXAppSecret + "&js_code=" + form.Code + "&grant_type=authorization_code")
	if err != nil {
		client.JSON(http.StatusBadRequest, gin.H{
			"code":    101,
			"message": err,
		})
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			client.JSON(http.StatusBadRequest, gin.H{
				"code":    106,
				"message": err,
			})
			return
		}
	}(resp.Body)

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		client.JSON(http.StatusBadRequest, gin.H{
			"code":    106,
			"message": err,
		})
		return
	}

	// 对resp转为JSON
	var p code2Session
	err = json.Unmarshal(b, &p)
	if err != nil {
		client.JSON(http.StatusBadRequest, gin.H{
			"code":    102,
			"message": err,
		})
		return
	}

	//设置登录信息缓存
	token := &Token{
		Session: p.SessionKey,
		Openid:  p.Openid,
	}

	cache := models.Cache{
		Key:      p.SessionKey,
		Value:    token,
		ExpireAt: time.Now().Add(time.Hour * 24 * 365),
	}

	if err = cache.Create(); err != nil {
		client.JSON(http.StatusBadRequest, gin.H{
			"code":    201,
			"message": "cache save error",
		})
	}

	client.JSON(http.StatusOK, gin.H{
		"code":    0,
		"openid":  p.Openid,
		"session": p.SessionKey,
	})
}
