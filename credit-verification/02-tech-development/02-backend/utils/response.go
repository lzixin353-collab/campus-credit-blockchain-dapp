package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构体
type Response struct {
	Code int         `json:"code"` // 状态码：200成功，其他失败
	Msg  string      `json:"msg"`  // 提示信息
	Data interface{} `json:"data"` // 业务数据（可选）
}

// 常用状态码定义
const (
	SuccessCode    = 200
	ErrorCode      = 500
	ParamErrorCode = 400
	AuthErrorCode  = 401
	ForbidCode     = 403
)

// Success 成功响应
func Success(c *gin.Context, data interface{}, msg ...string) {
	resMsg := "操作成功"
	if len(msg) > 0 {
		resMsg = msg[0]
	}
	c.JSON(http.StatusOK, Response{
		Code: SuccessCode,
		Msg:  resMsg,
		Data: data,
	})
}

// Error 失败响应
func Error(c *gin.Context, msg ...string) {
	resMsg := "操作失败"
	if len(msg) > 0 {
		resMsg = msg[0]
	}
	c.JSON(http.StatusOK, Response{
		Code: ErrorCode,
		Msg:  resMsg,
		Data: nil,
	})
}

// ParamError 参数错误响应
func ParamError(c *gin.Context, msg ...string) {
	resMsg := "参数错误"
	if len(msg) > 0 {
		resMsg = msg[0]
	}
	c.JSON(http.StatusOK, Response{
		Code: ParamErrorCode,
		Msg:  resMsg,
		Data: nil,
	})
}

// AuthError 认证失败响应（token无效/过期）
func AuthError(c *gin.Context, msg ...string) {
	resMsg := "登录失效，请重新登录"
	if len(msg) > 0 {
		resMsg = msg[0]
	}
	c.JSON(http.StatusOK, Response{
		Code: AuthErrorCode,
		Msg:  resMsg,
		Data: nil,
	})
}

// ForbidError 权限不足响应
func ForbidError(c *gin.Context, msg ...string) {
	resMsg := "权限不足，无法操作"
	if len(msg) > 0 {
		resMsg = msg[0]
	}
	c.JSON(http.StatusOK, Response{
		Code: ForbidCode,
		Msg:  resMsg,
		Data: nil,
	})
}
