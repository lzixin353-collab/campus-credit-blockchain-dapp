// utils/response.go
package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构体
type Response struct {
	Code int         `json:"code"` // 状态码（200成功，其他失败）
	Msg  string      `json:"msg"`  // 提示信息
	Data interface{} `json:"data"` // 响应数据
}

// Success 成功响应
func Success(c *gin.Context, data interface{}, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  msg,
		Data: data,
	})
}

// Fail 失败响应
func Fail(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: 400, // 通用失败码，可根据场景调整为401/403/500
		Msg:  msg,
		Data: nil,
	})
}

// FailWithCode 自定义状态码的失败响应（可选，补充用）
func FailWithCode(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}
