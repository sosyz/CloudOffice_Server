package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"sonui.cn/cloudprint/utils"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// DB 数据库链接单例
var DB *gorm.DB

// Init 初始化 MySQL 链接
func init() {

	var (
		db  *gorm.DB
		err error
	)
	db, err = gorm.Open(utils.Config.Run.DatabaseType, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		utils.Config.Db.User,
		utils.Config.Db.Password,
		utils.Config.Db.Host,
		utils.Config.Db.Port,
		utils.Config.Db.Database))
	if err != nil {
		panic(fmt.Sprintf("models.init err: %v", err))
	}

	//设置连接池
	db.DB().SetMaxOpenConns(100)

	//超时
	db.DB().SetConnMaxLifetime(time.Second * 30)
	// 设置表前缀
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "print_" + defaultTableName
	}
	// 更新结构
	db.AutoMigrate(&User{}, &File{}, &Auth{}, &Cache{}, &Order{})
	DB = db
}
