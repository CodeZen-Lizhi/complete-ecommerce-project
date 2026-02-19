package handler

import (
	"ecommerce/internal/response"

	"github.com/gin-gonic/gin"
)

// Result 定义统一响应结构。
type Result = response.Result

// Success 返回成功响应。
func Success(c *gin.Context, data interface{}) {
	response.Success(c, data)
}

// SuccessMsg 返回带自定义文案的成功响应。
func SuccessMsg(c *gin.Context, msg string, data interface{}) {
	response.SuccessMsg(c, msg, data)
}

// Fail 返回通用失败响应。
func Fail(c *gin.Context, msg string) {
	response.Fail(c, msg)
}

// Error 返回系统错误响应。
func Error(c *gin.Context, msg string) {
	response.Error(c, msg)
}

// ParamError 返回参数校验错误响应。
func ParamError(c *gin.Context, msg string) {
	response.ParamError(c, msg)
}

// AuthError 返回鉴权失败响应。
func AuthError(c *gin.Context, msg string) {
	response.AuthError(c, msg)
}

// ForbiddenError 返回权限不足响应。
func ForbiddenError(c *gin.Context, msg string) {
	response.ForbiddenError(c, msg)
}

// SystemError 返回包含请求链路ID的系统错误响应。
func SystemError(c *gin.Context, msg string, requestID string) {
	response.SystemError(c, msg, requestID)
}
