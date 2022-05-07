package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"sonui.cn/cloudprint/models"
	"sonui.cn/cloudprint/pkg/conf"
	"sonui.cn/cloudprint/pkg/utils"
	"sonui.cn/cloudprint/routers"
)

func main() {
	// TODO: 对配置类进行重构，使用VIPER库
	// 初始化配置
	conf.Type = "YAML"
	if err := conf.InitConfig("config.yaml"); err != nil {
		panic(fmt.Sprintf("init config failed, err: %s", err.Error()))
	}

	// 初始化雪花ID生成器
	if err := utils.NewWorker(conf.Conf.Config.Node); err != nil {
		panic(fmt.Sprintf("init snowflake failed, err: %s", err.Error()))
	}

	//初始化数据库连接
	models.Init()

	// 初始化gin
	gin.SetMode(gin.ReleaseMode)
	api := routers.InitRouter()

	// 启动服务
	if err := api.Run(conf.Conf.Config.Listen); err != nil {
		panic(err)
	}
}
