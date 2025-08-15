package router

import (
	"ecommerce/handler"
	"github.com/gin-gonic/gin"
)

// 注册订单模块的私有路由（需登录）
func registerOrderPrivateRoutes(private *gin.RouterGroup) {
	orders := private.Group("/orders")
	{
		orders.POST("", handler.CreateOrder)            // 创建订单
		orders.GET("", handler.ListUserOrders)          // 订单列表
		orders.GET("/:id", handler.GetOrderDetail)      // 订单详情
		orders.POST("/:id/cancel", handler.CancelOrder) // 取消订单
		orders.POST("/:id/pay", handler.PayOrder)       // 支付订单
	}
}

// 注册订单模块的管理员路由
func registerOrderAdminRoutes(admin *gin.RouterGroup) {
	orders := admin.Group("/orders")
	{
		orders.GET("", handler.ListAllOrders)                // 所有订单列表
		orders.PUT("/:id/status", handler.UpdateOrderStatus) // 更新订单状态
		orders.GET("/statistics", handler.OrderStatistics)   // 订单统计
	}
}
