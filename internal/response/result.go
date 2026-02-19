package response

import (
	"ecommerce/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Result 定义统一响应结构。
type Result struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	Data      interface{} `json:"data"`
	RequestID string      `json:"request_id,omitempty"`
}

// Success 返回默认成功响应。
func Success(c *gin.Context, data interface{}) {
	writeResult(c, http.StatusOK, 200, "success", data, "", false)
}

// SuccessMsg 返回自定义成功文案响应。
func SuccessMsg(c *gin.Context, msg string, data interface{}) {
	writeResult(c, http.StatusOK, 200, msg, data, "", false)
}

// Fail 返回通用失败响应。
func Fail(c *gin.Context, msg string) {
	writeResult(c, http.StatusOK, 400, msg, nil, "", false)
}

// Error 返回系统错误响应。
func Error(c *gin.Context, msg string) {
	writeResult(c, http.StatusOK, 500, msg, nil, "", false)
}

// ParamError 返回参数校验错误响应。
func ParamError(c *gin.Context, msg string) {
	writeResult(c, http.StatusBadRequest, 400, msg, nil, "", true)
}

// AuthError 返回鉴权失败响应。
func AuthError(c *gin.Context, msg string) {
	writeResult(c, http.StatusUnauthorized, 401, msg, nil, "", true)
}

// ForbiddenError 返回权限不足响应。
func ForbiddenError(c *gin.Context, msg string) {
	writeResult(c, http.StatusForbidden, 403, msg, nil, "", true)
}

// SystemError 返回包含请求链路ID的系统错误响应。
func SystemError(c *gin.Context, msg string, requestID string) {
	writeResult(c, http.StatusInternalServerError, 500, msg, nil, requestID, true)
}

// writeResult 统一写入响应体和 HTTP 状态码。
func writeResult(c *gin.Context, httpStatus int, code int, msg string, data interface{}, requestID string, abort bool) {
	r := &Result{
		Code:      code,
		Msg:       msg,
		Data:      data,
		RequestID: requestID,
	}
	c.Set(util.ResponseCodeKey, code)
	if abort {
		c.AbortWithStatusJSON(httpStatus, r)
		return
	}
	c.JSON(httpStatus, r)
}
