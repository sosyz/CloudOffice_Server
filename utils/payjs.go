package utils

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"sort"
	"strings"
)

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
