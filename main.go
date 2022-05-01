package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"sonui.cn/cloudprint/models"
	"sonui.cn/cloudprint/pkg/conf"
	"sonui.cn/cloudprint/routers"
)

func main() {
	// 初始化配置
	err := conf.InitConfig("config.yaml")
	if err != nil {
		fmt.Println("config init error:", err)
		return
	}

	// 初始化数据库连接
	models.Init()

	// 初始化gin
	gin.SetMode(gin.ReleaseMode)
	api := routers.InitRouter()

	// 启动服务
	if err := api.Run(":9000"); err != nil {
		panic(err)
	}
}
