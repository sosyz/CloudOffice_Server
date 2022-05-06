package utils

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"os"
	"sonui.cn/cloudprint/pkg/conf"
	"sort"
	"strings"
)

type Tools struct {
}

// GetEnvDefault 获取系统环境变量
func GetEnvDefault(key, defVal string) string {
	val, ex := os.LookupEnv(key)
	if !ex {
		return defVal
	}
	return val
}

// PayJsSign PayJs签名算法
func PayJsSign(order map[string]string, key string) (sign string) {
	data := url.Values{}
	for k, v := range order {
		data.Add(k, v)
	}
	keys := make([]string, 0, 0)
	for key := range data {
		if data.Get(key) != "" {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	body := data.Encode()
	d, _ := url.QueryUnescape(body)
	d += "&key=" + key
	h := md5.New()
	h.Write([]byte(d))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

// GetFilePagesNum 通过腾讯云万象能力获取文件页数
func GetFilePagesNum(path string) (string, error) {
	u, _ := url.Parse("https://" + conf.Conf.Cos.Bucket + ".cos." + conf.Conf.Cos.Region + ".myqcloud.com")
	cu, _ := url.Parse("https://" + conf.Conf.Cos.Bucket + ".ci." + conf.Conf.Cos.Region + ".myqcloud.com")
	b := &cos.BaseURL{BucketURL: u, CIURL: cu}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  conf.Conf.Secret.SecretId,
			SecretKey: conf.Conf.Secret.SecretKey,
		},
	})
	opt := &cos.DocPreviewOptions{}
	resp, err := c.CI.DocPreview(context.Background(), path, opt)
	if err != nil {
		return "", err
	}
	num := resp.Header.Get("X-Total-Page")
	if num == "" {
		num = "0"
	}
	return num, nil
}
