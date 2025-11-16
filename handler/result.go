package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(c *gin.Context, data interface{}) {
	r := &Result{
		Code: 200,
		Msg:  "success",
		Data: data,
	}
	c.JSON(http.StatusOK, r)
}
func Fail(c *gin.Context, msg string) {
	r := &Result{
		Code: 400,
		Msg:  msg,
		Data: nil,
	}
	c.JSON(http.StatusOK, r)
}

func Error(c *gin.Context, msg string) {
	r := &Result{
		Code: 500,
		Msg:  msg,
		Data: nil,
	}
	c.JSON(http.StatusOK, r)
}
