package payjs

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// Sign PayJs签名算法
func Sign(order map[string]string, key string) (sign string) {
	keys := make([]string, 0, 0)
	for key := range order {
		if order[key] != "" && strings.ToLower(key) != "sign" {
			keys = append(keys, key+"="+strings.TrimSpace(order[key]))
		}
	}
	sort.Strings(keys)
	d := strings.Join(keys, "&")
	d += "&key=" + key

	md5bs := md5.Sum([]byte(d))
	md5res := hex.EncodeToString(md5bs[:])
	ret := strings.ToUpper(md5res)
	return ret
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
	fmt.Printf("CheckSign[data]: %+v\n", data)
	if signStr == "" {
		return false, fmt.Errorf("sign is empty")
	} else {
		return Sign(data, key) == signStr, nil
	}

}
