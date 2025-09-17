package router

import (
	"ecommerce/middleware"
	"github.com/gin-gonic/gin"
)

// InitTotalRouter 初始化总路由
func InitTotalRouter() *gin.Engine {
	r := gin.New()
	// 全局中间件：跨域处理
	r.Use(middleware.CorsMiddleware(), middleware.RequestID(), gin.Recovery(), middleware.GinLogger())

	// 1. 公共路由组（无需登录）
	public := r.Group("/api/public")
	{
		// 注册各模块的公共路由
		registerUserPublicRoutes(public)
		registerProductPublicRoutes(public)
	}

	// 2. 私有路由组（需登录）
	private := r.Group("/api/private")
	private.Use(middleware.AuthMiddleware()) // 添加认证中间件
	{
		registerUserPrivateRoutes(private)
		registerOrderPrivateRoutes(private)
	}

	// 3. 管理员路由组（需管理员权限）
	admin := r.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminAuthMiddleware()) // 双重验证
	{
		registerProductAdminRoutes(admin)
		registerOrderAdminRoutes(admin)
	}

	return r
}
