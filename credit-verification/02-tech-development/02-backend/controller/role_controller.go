// controller/role_controller.go
package controller

import (
	"campus-credit-backend/utils"

	"github.com/gin-gonic/gin"
)

// AssignRole 分配合约角色
func AssignRole(c *gin.Context) {
	var req struct {
		UserAddress string `json:"user_address" binding:"required"`
		Role        string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, "参数错误: "+err.Error())
		return
	}

	// 调用合约AssignRole方法
	txHash, err := utils.AssignRole(req.UserAddress, req.Role)
	if err != nil {
		// 服务器内部错误用500码
		utils.FailWithCode(c, 500, "分配角色失败: "+err.Error())
		return
	}
	utils.Success(c, gin.H{"tx_hash": txHash}, "角色分配成功")
}

// GetRole 查询合约角色
func GetRole(c *gin.Context) {
	address := c.Query("user_address")
	if address == "" {
		utils.Fail(c, "地址不能为空")
		return
	}

	role, err := utils.GetRole(address)
	if err != nil {
		utils.FailWithCode(c, 500, "查询角色失败: "+err.Error())
		return
	}
	utils.Success(c, gin.H{"role": role}, "查询成功")
}
