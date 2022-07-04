package utils

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/viper"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var Config *config
var FileSF *Worker
var OrderSF *Worker
var S3 *session.Session

func init() {
	// 初始化配置文件
	Config = &config{}
	err := ReadConfig()
	if err != nil {
		fmt.Println("Read config failed, err:", err)
	}

	// 初始化雪花
	// string转为int64
	node, _ := strconv.ParseInt(Config.Run.Node, 10, 64)
	FileSF, _ = NewWorker(node)
	OrderSF, _ = NewWorker(node)

	if Config.Run.Temp == "" {
		Config.Run.Temp = "./temp"
	}

	if Config.Run.Temp[len(Config.Run.Temp)-1] != '/' {
		Config.Run.Temp += "/"
	}

	creds := credentials.NewStaticCredentials(Config.QCloud.SecretId, Config.QCloud.SecretKey, "")
	endpoint := "http://cos." + Config.QCloud.Region + ".myqcloud.com"
	conf := &aws.Config{
		Region:           aws.String(Config.QCloud.Region),
		Endpoint:         &endpoint,
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      creds,
		// DisableSSL:       &disableSSL,
	}
	S3, err = session.NewSession(conf)
	if err != nil {
		fmt.Println("Init s3 session failed, err:", err)
	}
}

func ReadConfig() error {
	//读取yaml文件
	v := viper.New()
	//设置读取的配置文件
	v.SetConfigName("config")
	//添加读取的配置文件路径
	v.AddConfigPath("./")
	//设置配置文件类型
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		// 从文件读取配置失败 尝试从环境变量获取
		fmt.Printf("Read config from env failed, err: %v, try to read from env\n", err)
		//设置环境变量名前缀
		return ReadConfigFromEnv(Config, "ONLINE_OFFICE")
	} else {
		if err := v.Unmarshal(Config); err != nil {
			return err
		}
	}

	return nil
}

// ReadConfigFromEnv 从环境变量读取配置
func ReadConfigFromEnv(obj interface{}, envPrefix string) error {
	objType := reflect.TypeOf(obj)
	objValue := reflect.ValueOf(obj)
	if objType.Kind() != reflect.Ptr {
		return fmt.Errorf("obj must be a pointer")
	}

	objType = objType.Elem()
	for i := 0; i < objType.NumField(); i++ {
		f := objType.Field(i)
		// 判断类型是否属于interface
		if f.Type.Kind() == reflect.Struct {
			// 如果是结构体，递归调用 传递指针
			if err := ReadConfigFromEnv(objValue.Elem().Field(i).Addr().Interface(), strings.ToUpper(envPrefix+"_"+f.Name)); err != nil {
				return err
			}
		} else {
			// 获取环境变量值
			val := os.Getenv(strings.ToUpper(envPrefix + "_" + f.Name))
			// 修改其值
			objValue.Elem().Field(i).SetString(val)
		}
	}
	return nil
}
