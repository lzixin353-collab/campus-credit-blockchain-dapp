// middleware/auth.go
package middleware

import (
	"net/http"
	"strings"

	"campus-credit-backend/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 登录鉴权中间件（验证JWT Token）
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header获取Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, utils.Response{
				Code: 401,
				Msg:  "未登录，请先登录",
				Data: nil,
			})
			c.Abort()
			return
		}

		// 解析Token格式（Bearer Token）
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, utils.Response{
				Code: 401,
				Msg:  "Token格式错误",
				Data: nil,
			})
			c.Abort()
			return
		}

		// 解析Token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.Response{
				Code: 401,
				Msg:  "Token无效或已过期",
				Data: nil,
			})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("userId", claims.UserId)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// RoleMiddleware 角色权限校验中间件
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, utils.Response{
				Code: 403,
				Msg:  "无权限访问",
				Data: nil,
			})
			c.Abort()
			return
		}

		// 校验角色是否在允许列表
		allow := false
		for _, r := range allowedRoles {
			if role.(string) == r {
				allow = true
				break
			}
		}
		if !allow {
			c.JSON(http.StatusForbidden, utils.Response{
				Code: 403,
				Msg:  "仅" + strings.Join(allowedRoles, "/") + "可访问",
				Data: nil,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
