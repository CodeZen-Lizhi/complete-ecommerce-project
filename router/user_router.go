package router

import (
	"ecommerce/handler"
	"github.com/gin-gonic/gin"
)

// 注册用户模块的公共路由（无需登录）
func registerUserPublicRoutes(public *gin.RouterGroup) {
	user := public.Group("/user")
	{
		user.POST("/register", handler.UserRegister) // 用户注册
		user.POST("/login", handler.UserLogin)       // 用户登录
		user.GET("/captcha", handler.GetCaptcha)     // 获取验证码
	}
}

// 注册用户模块的私有路由（需登录）
func registerUserPrivateRoutes(private *gin.RouterGroup) {
	user := private.Group("/user")
	{
		user.GET("/info", handler.GetUserInfo)     // 获取个人信息
		user.PUT("/info", handler.UpdateUserInfo)  // 更新个人信息
		user.PUT("/password", handler.ChangePwd)   // 修改密码
		user.GET("/orders", handler.UserOrderList) // 查看我的订单
	}
}
