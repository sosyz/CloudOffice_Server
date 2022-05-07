package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"sonui.cn/cloudprint/pkg/conf"
)

// DB 数据库链接单例
var DB *gorm.DB

// Init 初始化 MySQL 链接
func Init() {

	var (
		db  *gorm.DB
		err error
	)

	db, err = gorm.Open(conf.Conf.Config.DatabaseType, fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Conf.Db.User,
		conf.Conf.Db.Password,
		conf.Conf.Db.Host,
		conf.Conf.Db.Port,
		conf.Conf.Db.Database))
	if err != nil {
		panic(fmt.Sprintf("models.init err: %v", err))
	}

	//设置连接池
	db.DB().SetMaxOpenConns(100)

	//超时
	db.DB().SetConnMaxLifetime(time.Second * 30)
	// 检测表并创建
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "print_" + defaultTableName
	}

	// 检测表并创建
	if !db.HasTable(&User{}) {
		db.CreateTable(&User{})
	}
	if !db.HasTable(&File{}) {
		db.CreateTable(&File{})
	}
	if !db.HasTable(&Auth{}) {
		db.CreateTable(&Auth{})
	}
	if !db.HasTable(&Cache{}) && conf.Conf.Config.CacheType == "db" {
		db.CreateTable(&Cache{})
	}
	if !db.HasTable(&Order{}) {
		db.CreateTable(&Order{})
	}

	// 更新结构
	db.AutoMigrate(&User{}, &File{}, &Auth{}, &Cache{}, &Order{})
	DB = db
}
