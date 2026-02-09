// main.go
package main

import (
	"log" // 补充导入log包（原代码中用到log.Printf）

	"campus-credit-backend/router"
	"campus-credit-backend/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 初始化配置、数据库、以太坊客户端
	utils.InitConfig()
	utils.InitMySQL()
	utils.InitEthClient() // 你的原有以太坊客户端初始化

	// 2. 设置Gin运行模式（核心修复：改为包级别的gin.SetMode）
	gin.SetMode(utils.GlobalConfig.Server.Mode) // 关键修正！

	// 3. 初始化Gin引擎
	r := gin.Default()

	// 4. 初始化路由
	router.InitRouter(r)

	// 5. 启动服务
	port := utils.GlobalConfig.Server.Port
	log.Printf("服务启动成功，监听端口: %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
