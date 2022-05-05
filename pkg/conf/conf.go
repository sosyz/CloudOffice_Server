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

type conf struct {
	CacheType    string `yaml:"cache_type"`
	DatabaseType string `yaml:"database_type"`
	Node         int64  `yaml:"node"`
	Port         string `yaml:"port"`
}

// 配置
type config struct {
	Config conf   `yaml:"conf"`
	Cos    cos    `yaml:"cos"`
	Db     db     `yaml:"db"`
	Redis  redis  `yaml:"redis"`
	Secret secret `yaml:"secret"`
	Wechat wechat `yaml:"wechat"`
}

var Conf *config
var Type string

func InitConfig(path string) error {
	Conf = &config{}
	if Type == "Yaml" {
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
	Conf.Db.Host = utils.GetEnvDefault("DB_HOST", "")
	Conf.Db.Port, _ = strconv.Atoi(utils.GetEnvDefault("DB_PORT", "3306"))
	Conf.Db.User = utils.GetEnvDefault("DB_USER", "")
	Conf.Db.Password = utils.GetEnvDefault("DB_PASSWORD", "")
	Conf.Db.Database = utils.GetEnvDefault("DB_DATABASE", "")
	Conf.Db.Dir = utils.GetEnvDefault("DB_DIR", "")
	Conf.Redis.Host = utils.GetEnvDefault("REDIS_HOST", "")
	Conf.Redis.Port, _ = strconv.Atoi(utils.GetEnvDefault("REDIS_PORT", "0"))
	Conf.Redis.Db, _ = strconv.Atoi(utils.GetEnvDefault("REDIS_DB", "0"))
	Conf.Redis.Password = utils.GetEnvDefault("REDIS_PASSWORD", "")
	Conf.Secret.SecretId = utils.GetEnvDefault("SECRET_ID", "")
	Conf.Secret.SecretKey = utils.GetEnvDefault("SECRET_KEY", "")
	Conf.Cos.Appid = utils.GetEnvDefault("COS_APPID", "")
	Conf.Cos.Bucket = utils.GetEnvDefault("COS_BUCKET", "")
	Conf.Cos.Region = utils.GetEnvDefault("COS_REGION", "")
	Conf.Wechat.Appid = utils.GetEnvDefault("WECHAT_APPID", "")
	Conf.Wechat.Secret = utils.GetEnvDefault("WECHAT_SECRET", "")
	Conf.Config.CacheType = utils.GetEnvDefault("CACHE_TYPE", "")
	Conf.Config.DatabaseType = utils.GetEnvDefault("DATABASE_TYPE", "")
	Conf.Config.Node, _ = strconv.ParseInt(utils.GetEnvDefault("NODE", "0"), 10, 64)
	Conf.Config.Port = utils.GetEnvDefault("RUN_PORT", "")
	return nil
}
