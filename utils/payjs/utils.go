package payjs

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strings"
)

// Sign PayJs签名算法
func Sign(order map[string]string, key string) (sign string) {
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

// CheckSign 签名验证
func CheckSign(obj interface{}, key string) (bool, error) {
	data := map[string]string{}
	objType := reflect.TypeOf(obj)
	objValue := reflect.ValueOf(obj)
	if objType.Kind() != reflect.Ptr {
		return false, fmt.Errorf("obj must be a pointer")
	}

	objType = objType.Elem()
	signStr := ""
	for i := 0; i < objType.NumField(); i++ {
		f := objType.Field(i)
		if f.Name != "Sign" {
			switch f.Type.Kind() {
			case reflect.String:
				data[f.Tag.Get("json")] = fmt.Sprintf("%s", objValue.Elem().Field(i).String())
				break
			case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
				data[f.Tag.Get("json")] = fmt.Sprintf("%d", objValue.Elem().Field(i).Int())
				break
			case reflect.Float32, reflect.Float64:
				data[f.Tag.Get("json")] = fmt.Sprintf("%f", objValue.Elem().Field(i).Float())
				break
			case reflect.Bool:
				data[f.Tag.Get("json")] = fmt.Sprintf("%t", objValue.Elem().Field(i).Bool())
				break
			}
			//data[f.Tag.Get("json")] = objValue.Elem().Field(i).String()
		} else {
			signStr = objValue.Elem().Field(i).String()
		}
	}
	if signStr == "" {
		return false, fmt.Errorf("sign is empty")
	} else {
		return Sign(data, key) == signStr, nil
	}

}
