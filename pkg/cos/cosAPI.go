package cos

import (
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"sonui.cn/cloudprint/pkg/conf"
	"strconv"
)

// GetFilePagesNum 通过腾讯云万象能力获取文件页数
func GetFilePagesNum(path string) (int, error) {
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
		return 0, err
	}
	// string到int
	var num int
	if num, err = strconv.Atoi(resp.Header.Get("X-Total-Page")); err != nil {
		return 0, err
	}
	return num, nil
}
