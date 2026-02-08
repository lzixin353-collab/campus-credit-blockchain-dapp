package middleware

import (
	"campus-credit-backend/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件（验证token有效性）
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token（格式：Bearer xxxx）
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.AuthError(c, "请先登录")
			c.Abort() // 终止请求
			return
		}

		// 解析token格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			utils.AuthError(c, "token格式错误")
			c.Abort()
			return
		}

		// 解析token
		tokenString := parts[1]
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			utils.AuthError(c, "token无效或已过期")
			c.Abort()
			return
		}

		// 将用户信息存入上下文，供后续接口使用
		c.Set("user_id", claims.UserID)
		c.Set("address", claims.Address)
		c.Set("role", claims.Role)

		c.Next() // 继续执行后续中间件/接口
	}
}

// RoleAuth 角色权限校验中间件（指定允许的角色）
func RoleAuth(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取用户角色
		role, exists := c.Get("role")
		if !exists {
			utils.ForbidError(c)
			c.Abort()
			return
		}

		// 校验角色是否在允许列表中
		allowed := false
		userRole := role.(string)
		for _, r := range roles {
			if userRole == r {
				allowed = true
				break
			}
		}

		if !allowed {
			utils.ForbidError(c)
			c.Abort()
			return
		}

		c.Next()
	}
}
