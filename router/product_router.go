package router

import (
	"ecommerce/handler"
	"github.com/gin-gonic/gin"
)

// 注册商品模块的公共路由（无需登录）
func registerProductPublicRoutes(public *gin.RouterGroup) {
	products := public.Group("/products")
	{
		products.GET("", handler.ListProducts)         // 商品列表
		products.GET("/:id", handler.GetProductDetail) // 商品详情
	}

	categories := public.Group("/categories")
	{
		categories.GET("", handler.ListCategories)                      // 分类列表
		categories.GET("/:id/products", handler.ListProductsByCategory) // 分类下的商品
	}
}

// 注册商品模块的管理员路由
func registerProductAdminRoutes(admin *gin.RouterGroup) {
	products := admin.Group("/products")
	{
		products.POST("", handler.CreateProduct)       // 创建商品
		products.PUT("/:id", handler.UpdateProduct)    // 更新商品
		products.DELETE("/:id", handler.DeleteProduct) // 删除商品
	}

	categories := admin.Group("/categories")
	{
		categories.POST("", handler.CreateCategory)       // 创建分类
		categories.PUT("/:id", handler.UpdateCategory)    // 更新分类
		categories.DELETE("/:id", handler.DeleteCategory) // 删除分类
	}
}
