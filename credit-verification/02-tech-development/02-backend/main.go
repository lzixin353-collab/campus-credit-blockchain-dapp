package main

import (
	"campus-credit-backend/middleware"
	"campus-credit-backend/utils"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 初始化配置
	utils.InitConfig()

	// 2. 初始化MySQL
	utils.InitMySQL()

	// 新增：初始化以太坊客户端
	utils.InitEthClient()

	// 3. 设置Gin模式
	gin.SetMode(utils.GlobalConfig.Server.Mode)

	// 4. 创建Gin引擎
	r := gin.Default()

	// 5. 配置跨域（原有代码不变）
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 6. 路由部分（原有代码不变）
	r.GET("/ping", func(c *gin.Context) {
		utils.Success(c, gin.H{"message": "pong"}, "接口访问成功")
	})

	// 7. 需要认证的路由组（原有代码不变）
	authGroup := r.Group("/api")
	// authGroup.Use(middleware.JWTAuth())
	{
		authGroup.GET("/user/info", func(c *gin.Context) {
			userID := c.GetUint64("user_id")
			address := c.GetString("address")
			role := c.GetString("role")
			utils.Success(c, gin.H{
				"user_id": userID,
				"address": address,
				"role":    role,
			}, "获取用户信息成功")
		})

		adminGroup := authGroup.Group("/admin")
		adminGroup.Use(middleware.RoleAuth("admin"))
		{
			adminGroup.GET("/user/list", func(c *gin.Context) {
				utils.Success(c, gin.H{"list": []string{}}, "管理员获取用户列表成功")
			})
		}
	}

	// 在authGroup后新增合约交互路由
	// 8. 合约交互路由（需要认证）
	contractGroup := authGroup.Group("/contract")
	{
		// 分配角色
		contractGroup.POST("/assign-role", func(c *gin.Context) {
			var req struct {
				UserAddress string `json:"user_address" binding:"required"`
				Role        string `json:"role" binding:"required,oneof=admin teacher student"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				utils.ParamError(c, err.Error())
				return
			}

			txHash, err := utils.AssignRole(req.UserAddress, req.Role)
			if err != nil {
				utils.Error(c, err.Error())
				return
			}

			utils.Success(c, gin.H{"tx_hash": txHash}, "角色分配成功")
		})

		// 录入学分
		contractGroup.POST("/record-credit", func(c *gin.Context) {
			var req struct {
				UserAddress string  `json:"user_address" binding:"required"`
				CourseName  string  `json:"course_name" binding:"required"`
				Score       float64 `json:"score" binding:"required,min=0,max=100"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				utils.ParamError(c, err.Error())
				return
			}

			txHash, err := utils.RecordCredit(req.UserAddress, req.CourseName, req.Score)
			if err != nil {
				utils.Error(c, err.Error())
				return
			}

			utils.Success(c, gin.H{"tx_hash": txHash}, "学分录入成功")
		})

		// 查询用户学分
		contractGroup.GET("/user-credits/:address", func(c *gin.Context) {
			address := c.Param("address")
			if address == "" {
				utils.ParamError(c, "用户地址不能为空")
				return
			}

			credits, err := utils.GetUserCredits(address)
			if err != nil {
				utils.Error(c, err.Error())
				return
			}

			utils.Success(c, gin.H{"credits": credits}, "查询学分成功")
		})

		// 审核学分
		contractGroup.POST("/audit-credit", func(c *gin.Context) {
			var req struct {
				CreditId uint64 `json:"credit_id" binding:"required"`
				Approved bool   `json:"approved" binding:"required"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				utils.ParamError(c, err.Error())
				return
			}

			txHash, err := utils.AuditCredit(req.CreditId, req.Approved)
			if err != nil {
				utils.Error(c, err.Error())
				return
			}

			utils.Success(c, gin.H{"tx_hash": txHash}, "学分审核成功")
		})
	}

	// 9. 启动服务器
	port := utils.GlobalConfig.Server.Port
	log.Printf("服务器启动在端口: %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
