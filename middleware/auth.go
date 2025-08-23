package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// CorsMiddleware 处理跨域请求
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许所有来源
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// 允许的请求方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// 允许的请求头
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, Role")
		// 允许暴露的响应头
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		// 允许携带cookie
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// 处理OPTIONS请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// AuthMiddleware 验证用户是否登录
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		token := c.GetHeader("Authorization")
		// 简单验证token（实际项目中应该解析和验证token）
		/*if token == "" || token != "valid_token" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "未授权访问，请先登录",
			})
			c.Abort()
			return
		}*/
		// 第二种写法
		if token != "valid_token" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "未授权访问，请先登录",
			})
			return
		}

		// 可以在这里解析token获取用户信息并设置到上下文
		c.Set("userID", 10001) // 示例用户ID
		c.Next()
	}
}

// AdminAuthMiddleware 验证是否为管理员
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Role头
		role := c.GetHeader("Role")

		// 验证是否为管理员
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "没有管理员权限，无法访问",
			})
			c.Abort()
			return
		}

		c.Next()
	}
	//todo 写一个统计方法用时的中间件
}
