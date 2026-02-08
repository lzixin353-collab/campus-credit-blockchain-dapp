package main

import (
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

	// 3. 设置Gin模式
	gin.SetMode(utils.GlobalConfig.Server.Mode)

	// 4. 创建Gin引擎
	r := gin.Default()

	// 5. 配置跨域
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 开发环境允许所有，生产需指定前端域名
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 6. 测试路由
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "pong",
			"data": nil,
		})
	})

	// 7. 启动服务器
	port := utils.GlobalConfig.Server.Port
	log.Printf("服务器启动在端口: %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
