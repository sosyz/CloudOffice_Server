package conf

import (
	"gopkg.in/yaml.v2"
	"sonui.cn/cloudprint/pkg/fileTools"
	"sonui.cn/cloudprint/pkg/utils"
	"strconv"
)

// cos配置
type cos struct {
	Appid  string `yaml:"appid"`
	Bucket string `yaml:"bucket"`
	Region string `yaml:"region"`
}

// 微信
type wechat struct {
	Appid  string `yaml:"appid"`
	Secret string `yaml:"secret"`
}

// 数据库配置
type db struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Dir      string `yaml:"dir"`
}

// Redis配置
type redis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Db       int    `yaml:"db"`
	Password string `yaml:"password"`
}

// 访问令牌配置
type secret struct {
	SecretId  string `yaml:"secret_id"`
	SecretKey string `yaml:"secret_key"`
}

// 配置
type conf struct {
	CacheType    string `yaml:"cache_type"`
	DatabaseType string `yaml:"database_type"`
	Node         int64  `yaml:"node"`
	Listen       string `yaml:"listen"`
}

// 支付配置
type pay struct {
	MchId string `yaml:"mch_id"`
	Key   string `yaml:"key"`
}

// 入口
type config struct {
	Config conf   `yaml:"conf"`
	Cos    cos    `yaml:"cos"`
	Db     db     `yaml:"db"`
	Redis  redis  `yaml:"redis"`
	Secret secret `yaml:"secret"`
	Wechat wechat `yaml:"wechat"`
	Pay    pay    `yaml:"pay"`
}

var Conf *config
var Type string

func InitConfig(path string) error {
	// TODO: 配置转为VIPER库
	Conf = &config{}
	if Type == "YAML" {
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
	} else if Type == "ENV" {
		return initFromEnv()
	}
	return nil
}

func initFromEnv() error {
	// TODO: redis配置读取
	Conf.Db.Host = utils.GetEnvDefault("db_host", "127.0.0.1")
	Conf.Db.Port, _ = strconv.Atoi(utils.GetEnvDefault("db_port", "3306"))
	Conf.Db.User = utils.GetEnvDefault("db_user", "")
	Conf.Db.Password = utils.GetEnvDefault("db_password", "")
	Conf.Db.Database = utils.GetEnvDefault("db_database", "")

	Conf.Secret.SecretId = utils.GetEnvDefault("secret_secret_id", "")
	Conf.Secret.SecretKey = utils.GetEnvDefault("secret_secret_key", "")

	Conf.Cos.Appid = utils.GetEnvDefault("cos_appid", "")
	Conf.Cos.Bucket = utils.GetEnvDefault("cos_bucket", "")
	Conf.Cos.Region = utils.GetEnvDefault("cos_region", "")

	Conf.Wechat.Appid = utils.GetEnvDefault("wechat_appid", "")
	Conf.Wechat.Secret = utils.GetEnvDefault("wechat_secret", "")

	Conf.Config.CacheType = utils.GetEnvDefault("conf_cache_type", "")
	Conf.Config.DatabaseType = utils.GetEnvDefault("conf_database_type", "")
	Conf.Config.Node, _ = strconv.ParseInt(utils.GetEnvDefault("conf_node", "0"), 10, 64)
	Conf.Config.Listen = utils.GetEnvDefault("conf_listen", "0.0.0.0:9000")

	Conf.Pay.MchId = utils.GetEnvDefault("pay_mch_id", "")
	Conf.Pay.Key = utils.GetEnvDefault("pay_key", "")

	return nil
}
