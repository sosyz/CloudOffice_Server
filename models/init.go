package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"sonui.cn/cloudprint/pkg/utils"
)

// DB 数据库链接单例
var DB *gorm.DB

// Init 初始化 MySQL 链接
func Init() {

	var (
		db  *gorm.DB
		err error
	)

	dbType := utils.GetEnvDefault("database", "DB_TYPE")
	if dbType == "UNSET" {
		db, err = gorm.Open("sqlite3", utils.GetEnvDefault("DBFile_PATH", "/tmp"))
	} else {
		db, err = gorm.Open(utils.GetEnvDefault("DB_TYPE", "mysql"), fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			utils.GetEnvDefault("DB_USER", "root"),
			utils.GetEnvDefault("DB_PASSWORD", "root"),
			utils.GetEnvDefault("DB_HOST", "127.0.0.1"),
			utils.GetEnvDefault("DB_NAME", "cloudprint"),
		))

	}
	if err != nil {
		panic(fmt.Sprintf("models.init err: %v", err))
	}

	//设置连接池
	db.DB().SetMaxIdleConns(50)
	if dbType == "sqlite" || dbType == "sqlite3" || dbType == "UNSET" {
		db.DB().SetMaxOpenConns(1)
	} else {
		db.DB().SetMaxOpenConns(100)
	}

	//超时
	db.DB().SetConnMaxLifetime(time.Second * 30)

	DB = db

	//TODO: 迁移数据库功能
	//migration()
}
