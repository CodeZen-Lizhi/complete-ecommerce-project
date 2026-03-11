package handler

import (
	"ecommerce/container"
	"ecommerce/internal/logger"
	"ecommerce/model"
	"ecommerce/service"
	"ecommerce/util"
	"errors"

	"github.com/gin-gonic/gin"
)

// UserRegister 用户注册
func UserRegister(c *gin.Context) {
	var userVO model.UserVo
	// 接收参数
	if err := c.ShouldBindJSON(&userVO); err != nil {
		ParamError(c, "参数错误")
		return
	}
	con, ok := getContainerOrFail(c)
	if !ok {
		return
	}
	log := logger.GetLogger()
	log.Info("用户注册", "username", userVO.Username, "email", userVO.Email, "age", userVO.Age)
	userService := con.UserService

	// 创建用户
	user := model.User{
		Username: userVO.Username,
		Email:    userVO.Email,
		Password: userVO.Password,
		Age:      userVO.Age,
	}
	if err := userService.Register(c.Request.Context(), &user); err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			Fail(c, "用户已存在")
			return
		}
		Error(c, "用户创建失败")
		return
	}
	SuccessMsg(c, "注册成功", gin.H{})
}

// UserLogin 用户登录
func UserLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	con, ok := getContainerOrFail(c)
	if !ok {
		return
	}
	user, err := con.UserService.Login(username, password)
	if err != nil {
		Fail(c, "用户名或密码错误")
		return
	}
	token, err := util.GenerateToken(user.ID)
	if err != nil {
		Fail(c, "生成token失败")
		return
	}
	Success(c, token)
}

// GetCaptcha 获取验证码
func GetCaptcha(c *gin.Context) {
	// 实际项目中会生成图片验证码
	Success(c, gin.H{
		"captchaId": "cap_123456",
		"imageUrl":  "https://example.com/captcha.jpg",
	})
}

// GetUserInfo 获取个人信息
func GetUserInfo(c *gin.Context) {
	log := logger.GetLogger()
	userID, exists := c.Get(util.CurrentUserId)
	if !exists {
		log.Error("未获取到用户 ID")
		Fail(c, "未获取到用户 ID")
		return
	}
	con, ok := getContainerOrFail(c)
	if !ok {
		return
	}
	// 获取用户服务
	userService := con.UserService
	user, err := userService.FindByID(userID.(int64))
	if err != nil {
		log.Error("用户不存在")
		Fail(c, "用户不存在")
		return
	}
	Success(c, user)
}

// getContainerOrFail 获取已注入容器，不存在时直接返回统一错误响应
func getContainerOrFail(c *gin.Context) (*container.Container, bool) {
	log := logger.GetLogger()
	con, err := container.GetContainer(c)
	if err != nil {
		log.Error("容器未注入", "error", err)
		Error(c, "系统繁忙，请稍后重试")
		return nil, false
	}
	return con, true
}

// UpdateUserInfo 更新个人信息
func UpdateUserInfo(c *gin.Context) {
	nickname := c.PostForm("nickname")
	avatar := c.PostForm("avatar")

	SuccessMsg(c, "信息更新成功", gin.H{
		"nickname": nickname,
		"avatar":   avatar,
	})
}

// ChangePwd 修改密码
func ChangePwd(c *gin.Context) {
	_ = c.PostForm("oldPassword")
	_ = c.PostForm("newPassword")

	// 实际项目中会验证旧密码并更新新密码
	SuccessMsg(c, "密码修改成功", gin.H{})
}

// UserOrderList 查看我的订单
func UserOrderList(c *gin.Context) {
	userID, _ := c.Get(util.CurrentUserId)

	Success(c, gin.H{
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
	})
}
