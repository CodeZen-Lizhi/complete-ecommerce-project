package handler

import (
	"ecommerce/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Result struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	Data      interface{} `json:"data"`
	RequestID string      `json:"request_id,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	writeResult(c, http.StatusOK, 200, "success", data, "", false)
}

func Fail(c *gin.Context, msg string) {
	writeResult(c, http.StatusOK, 400, msg, nil, "", false)
}

func Error(c *gin.Context, msg string) {
	writeResult(c, http.StatusOK, 500, msg, nil, "", false)
}

func ParamError(c *gin.Context, msg string) {
	writeResult(c, http.StatusBadRequest, 400, msg, nil, "", true)
}

func AuthError(c *gin.Context, msg string) {
	writeResult(c, http.StatusUnauthorized, 401, msg, nil, "", true)
}

func ForbiddenError(c *gin.Context, msg string) {
	writeResult(c, http.StatusForbidden, 403, msg, nil, "", true)
}

func SystemError(c *gin.Context, msg string, requestID string) {
	writeResult(c, http.StatusInternalServerError, 500, msg, nil, requestID, true)
}

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
