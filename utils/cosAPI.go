package utils

import (
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"strconv"
)

// GetFilePagesNum 通过腾讯云万象能力获取文件页数
func GetFilePagesNum(path string) (int, error) {
	u, _ := url.Parse("https://" + Config.QCloud.Bucket + ".cos." + Config.QCloud.Region + ".myqcloud.com")
	cu, _ := url.Parse("https://" + Config.QCloud.Bucket + ".ci." + Config.QCloud.Region + ".myqcloud.com")
	b := &cos.BaseURL{BucketURL: u, CIURL: cu}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  Config.QCloud.SecretId,
			SecretKey: Config.QCloud.SecretKey,
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
