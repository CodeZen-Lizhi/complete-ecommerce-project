package handler

import (
	"ecommerce/container"
	"ecommerce/internal/logger"
	"ecommerce/model"
	"ecommerce/repository"
	"ecommerce/service"
	"ecommerce/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserRegister 用户注册
func UserRegister(c *gin.Context) {
	var userVo model.UserVo
	//接收参数
	if err := c.ShouldBindJSON(&userVo); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "参数错误",
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}
	log := logger.GetLogger()
	log.Info("用户注册", "username", userVo.Username, "email", userVo.Email, "age", userVo.Age)
	dao := repository.NewUserRepository()
	service := service.NewUserService(dao)
	//ctn, err := GetContainer(c)

	//验证登录名是否存在
	if exists, err := service.IzExist(userVo.Username); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "用户查询失败",
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	} else if exists {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "用户已存在",
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}
	//新增 用户
	user := model.User{
		Username: userVo.Username,
		Email:    userVo.Email,
		Password: userVo.Password,
		Age:      userVo.Age,
	}
	if err := service.Create(&user); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "用户创建失败",
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "注册成功",
		"data": gin.H{},
	})
}

// UserLogin 用户登录
func UserLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	dao := repository.NewUserRepository()
	service := service.NewUserService(dao)
	user, err := service.Login(username, password)
	if err != nil {
		Fail(c, "用户名或密码错误")
		return
	}
	token, err := util.GenerateToken(strconv.FormatInt(user.ID, 10))
	if err != nil {
		Fail(c, "生成token失败")
		return
	}
	Success(c, token)
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
	// 从上下文获取用户ID
	log := logger.GetLogger()
	userID, exists := c.Get(util.CurrentUserId)
	if !exists {
		log.Error("未获取到用户id")
		Fail(c, "未获取到用户id")
		return
	}
	con, err2 := container.GetContainer(c)
	if err2 != nil {
		log.Error("容器未注入")
		Fail(c, "容器未注入")
		return
	}
	//手写注入
	userService := con.UserService
	user, err := userService.FindByID(userID.(uint64))
	if err != nil {
		log.Error("用户不存在")
		Fail(c, "用户不存在")
		return
	}
	Success(c, user)
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
	userID, _ := c.Get(util.CurrentUserId)

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
