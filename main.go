package main

import (
	"github.com/gin-gonic/gin"
	"sonui.cn/cloudprint/routers"
	"sonui.cn/cloudprint/utils"
)

func main() {
	// 初始化gin
	gin.SetMode(gin.ReleaseMode)
	api := routers.InitRouter()
	// 启动服务
	if err := api.Run(utils.Config.Run.Listen); err != nil {
		panic(err)
	}
}
