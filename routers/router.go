package routers

import (
	"github.com/gin-gonic/gin"
	"sonui.cn/cloudprint/routers/controllers"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	r := gin.Default()
	// r.Use(Cors())

	//服务器状态检测
	r.GET("/ping", func(context *gin.Context) {
		context.String(200, "!pong")
	})

	// v1接口
	v1 := r.Group("/api/v1")
	// Session校验
	v1.Use(controllers.TokenHandler())
	// 路由
	{
		// 用户模块
		user := v1.Group("/user")
		{
			// 登录
			login := user.Group("/login")
			{
				// 默认登录方式 使用openidId
				login.GET("default", controllers.LoginDefault)

				// 扫码登录 获取QR码
				login.GET("qr", controllers.LoginQR)

				// 扫码登录 允许登录
				login.GET("check", controllers.LoginCheck)

				// 扫码登录 检查状态
				login.GET("status", controllers.LoginStatus)
			}

			// 获取账号信息
			user.GET("info", controllers.UserGetInfo)

			// 设置账号信息
			user.POST("set", controllers.UserSetInfo)
		}

		file := v1.Group("/file")
		{
			// 上传文件
			file.POST("upload", controllers.FileUpload)

			// 下载文件
			file.GET("download", controllers.FileDownload)

			file.HEAD("download", controllers.FileDownload)

		}

		order := v1.Group("/order")
		{
			// 获取订单列表
			order.GET("list", controllers.OrderList)

			// 获取订单详情
			order.POST("detail", controllers.OrderDetail)

			// 合成订单
			order.POST("merge", controllers.OrderMerge)

			// 取消订单
			order.POST("cancel", controllers.OrderCancel)

			// 支付
			pay := order.Group("/pay")
			{
				// 获取订单支付信息
				pay.GET("info", controllers.OrderPayInfo)

				// 获取订单支付状态
				pay.GET("status", controllers.OrderPayStatus)

				// 支付回调
				pay.POST("notify", controllers.OrderPayNotify)
			}

			// 授权商家再次读取订单文件
			order.POST("repeatRead", controllers.OrderFileRepeatRead)
		}
	}
	return r

}
