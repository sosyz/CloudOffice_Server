package routers

import (
	"github.com/gin-gonic/gin"
	"sonui.cn/cloudprint/routers/controllers"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	r := gin.Default()
	// r.Use(Cors())

	// Session校验
	r.Use(controllers.TokenHandler())
	// v1接口
	v1 := r.Group("/api/v1")

	// 路由
	{
		// 用户模块
		user := v1.Group("/user")
		{
			// 登录
			user.POST("login", controllers.Login)

			// 获取临时密钥
			user.POST("tmpKey", controllers.CreatTmpKey)
		}

		file := v1.Group("/file")
		{
			// 通知开始上传文件
			file.POST("upload", controllers.UploadStart)

			// 通知文件上传完毕
			file.POST("complete", controllers.UploadComplete)
		}

		order := v1.Group("/order")
		{
			// 获取订单列表
			order.POST("list", controllers.OrderList)

			// 获取订单详情
			order.POST("detail", controllers.OrderDetail)

			// 合成订单
			order.POST("synthesis", controllers.OrderSynthesis)

			// 取消订单
			order.POST("cancel", controllers.OrderCancel)

			// 获取订单支付信息
			order.GET("payInfo", controllers.OrderPayInfo)

			// 获取订单支付状态
			order.GET("payStatus", controllers.OrderPayStatus)
		}

		ws := v1.Group("/wss")
		{
			// 打印机监听
			ws.GET("printer", controllers.Printer)

			// 用户监听
			ws.GET("user", controllers.User)
		}
	}
	return r

}
