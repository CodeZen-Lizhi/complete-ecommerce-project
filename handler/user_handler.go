package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// UserRegister 用户注册
func UserRegister(c *gin.Context) {
	username := c.PostForm("username")
	_ = c.PostForm("password")
	email := c.PostForm("email")

	// 实际项目中会有参数验证、数据库操作等
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "注册成功",
		"data": gin.H{
			"username": username,
			"email":    email,
			"createAt": time.Now().Format("2006-01-02 15:04:05"),
		},
	})
}

// UserLogin 用户登录
func UserLogin(c *gin.Context) {
	username := c.PostForm("username")
	_ = c.PostForm("password")

	// 实际项目中会验证用户名密码并生成token
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"token":    "valid_token", // 实际项目中是加密的token
			"username": username,
			"expireAt": time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
		},
	})
}

// GetCaptcha 获取验证码
func GetCaptcha(c *gin.Context) {
	// 实际项目中会生成图片验证码
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"captchaId": "cap_123456",
			"imageUrl":  "https://example.com/captcha.jpg",
		},
	})
}

// GetUserInfo 获取个人信息
func GetUserInfo(c *gin.Context) {
	// 从上下文获取用户ID（实际项目中从token解析）
	userID, _ := c.Get("userID")

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"id":       userID,
			"username": "testuser",
			"email":    "test@example.com",
			"phone":    "13800138000",
			"avatar":   "https://example.com/avatar.jpg",
			"joinTime": "2023-01-15",
		},
	})
}

// UpdateUserInfo 更新个人信息
func UpdateUserInfo(c *gin.Context) {
	nickname := c.PostForm("nickname")
	avatar := c.PostForm("avatar")

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "信息更新成功",
		"data": gin.H{
			"nickname": nickname,
			"avatar":   avatar,
		},
	})
}

// ChangePwd 修改密码
func ChangePwd(c *gin.Context) {
	_ = c.PostForm("oldPassword")
	_ = c.PostForm("newPassword")

	// 实际项目中会验证旧密码并更新新密码
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "密码修改成功",
	})
}

// UserOrderList 查看我的订单
func UserOrderList(c *gin.Context) {
	// 从上下文获取用户ID
	userID, _ := c.Get("userID")

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"userId": userID,
			"list": []gin.H{
				{
					"id":      1001,
					"orderNo": "ORD20230601001",
					"amount":  7999.00,
					"status":  "已支付",
					"time":    "2023-06-01 10:30:00",
				},
				{
					"id":      1002,
					"orderNo": "ORD20230610002",
					"amount":  1299.00,
					"status":  "已发货",
					"time":    "2023-06-10 15:20:00",
				},
			},
		},
	})
}
