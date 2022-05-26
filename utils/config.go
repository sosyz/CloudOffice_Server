package utils

// 腾讯云配置
type qcloud struct {
	Appid     string
	Bucket    string
	Region    string
	SecretId  string
	SecretKey string
}

// 微信
type wechat struct {
	Appid  string
	Secret string
}

// 数据库配置
type db struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Dir      string
}

// Redis配置
type redis struct {
	Host     string
	Port     string
	Db       string
	Password string
}

// 配置
type run struct {
	CacheType    string
	DatabaseType string
	Node         string
	Listen       string
}

// 支付配置
type pay struct {
	MchId string
	Key   string
}

// 入口
type config struct {
	Run    run
	QCloud qcloud
	Db     db
	Redis  redis
	Wechat wechat
	Pay    pay
}
