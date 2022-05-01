package conf

import (
	"gopkg.in/yaml.v2"
	"sonui.cn/cloudprint/pkg/fileTools"
)

// cos配置
type cos struct {
	Appid  string `yaml:"appid"`
	Bucket string `yaml:"bucket"`
	Region string `yaml:"region"`
}

// MySQL配置
type mysql struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

// Redis配置
type redis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	db       int    `yaml:"db"`
	Password string `yaml:"password"`
}

// 访问令牌配置
type secret struct {
	SecretId  string `yaml:"secret_id"`
	SecretKey string `yaml:"secret_key"`
}

type conf struct {
	CacheType string `yaml:"cache_type"`
}

// 配置
type config struct {
	Conf   conf   `yaml:"conf"`
	Cos    cos    `yaml:"cos"`
	Mysql  mysql  `yaml:"mysql"`
	Redis  redis  `yaml:"redis"`
	Secret secret `yaml:"secret"`
}

var Conf *config

func InitConfig(path string) error {
	Conf = &config{}
	var err error
	var data []byte

	if path == "" || !fileTools.Exists(path) {
		panic("config file not found")
	}

	data, err = fileTools.Read(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, Conf)
	if err != nil {
		return err
	}
	return nil
}
