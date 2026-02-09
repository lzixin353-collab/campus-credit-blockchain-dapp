// router/router.go
package router

import (
	"campus-credit-backend/controller"
	"campus-credit-backend/middleware"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter(r *gin.Engine) {
	// 公开接口（无需登录）
	public := r.Group("/api")
	{
		// 用户注册/登录
		public.POST("/user/register", controller.UserRegister)
		public.POST("/user/login", controller.UserLogin)
	}

	// 需登录的接口（全局鉴权）
	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())
	{
		// 用户管理
		user := auth.Group("/user")
		{
			user.GET("/info", controller.UserInfo)
			user.POST("/update", controller.UserUpdate)
		}

		// 角色管理（仅admin/teacher可访问）
		role := auth.Group("/role")
		role.Use(middleware.RoleMiddleware("admin", "teacher"))
		{
			role.POST("/assign", controller.AssignRole)
			role.GET("/get", controller.GetRole)
		}
	}
}
