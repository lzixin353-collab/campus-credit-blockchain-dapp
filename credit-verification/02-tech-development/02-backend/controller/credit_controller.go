// controller/credit_controller.go 学分录入、审核、查询、同步
package controller

import (
	"context"
	"time"

	"campus-credit-backend/model"
	"campus-credit-backend/utils"

	"github.com/gin-gonic/gin"
)

// CreditRecordReq 录入学分请求（教师）
type CreditRecordReq struct {
	StudentAddress string  `json:"student_address" binding:"required"`
	CourseName     string  `json:"course_name" binding:"required"`
	Score          float64 `json:"score" binding:"required,gte=0,lte=100"`
}

// CreditRecord 教师录入学分（上链 + 落库）
func CreditRecord(c *gin.Context) {
	userId, _ := c.Get("userId")
	user, err := model.GetUserById(userId.(uint64))
	if err != nil || user == nil {
		utils.Fail(c, "用户不存在")
		return
	}
	if !user.Address.Valid || user.Address.String == "" {
		utils.Fail(c, "请先绑定钱包地址")
		return
	}
	teacherAddress := user.Address.String

	var req CreditRecordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, "参数错误: "+err.Error())
		return
	}

	txHash, err := utils.RecordCredit(req.StudentAddress, req.CourseName, req.Score)
	if err != nil {
		utils.Fail(c, "上链失败: "+err.Error())
		return
	}

	// 等待交易打包后再查链上，否则 getStudentCredits 可能还读不到新数据
	if err := utils.WaitTxMined(context.Background(), txHash, 15*time.Second); err != nil {
		utils.Fail(c, "上链成功但等待打包超时，请稍后在「录入列表」查看")
		return
	}

	credits, err := utils.GetUserCredits(req.StudentAddress)
	if err != nil || len(credits) == 0 {
		utils.Fail(c, "上链成功但获取链上学分ID失败，请稍后同步")
		return
	}
	var maxId uint64
	for _, cr := range credits {
		if id, ok := cr["id"].(uint64); ok && id > maxId {
			maxId = id
		}
	}
	contractCreditId := int64(maxId)

	_, err = model.CreateCredit(req.StudentAddress, teacherAddress, req.CourseName, req.Score, "pending", txHash, contractCreditId)
	if err != nil {
		utils.Fail(c, "保存记录失败: "+err.Error())
		return
	}
	utils.Success(c, gin.H{"tx_hash": txHash, "contract_credit_id": contractCreditId}, "录入学分成功")
}

// CreditApproveReq 审核请求（管理员）
type CreditApproveReq struct {
	CreditId int64 `json:"credit_id" binding:"required"` // 数据库主键 id
}

// CreditApprove 管理员审核学分（调合约 + 更新库）
func CreditApprove(c *gin.Context) {
	var req CreditApproveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, "参数错误: "+err.Error())
		return
	}
	row, err := model.GetCreditById(req.CreditId)
	if err != nil || row == nil {
		utils.Fail(c, "学分记录不存在")
		return
	}
	if row.Status != "pending" {
		utils.Fail(c, "该记录已审核")
		return
	}
	contractId := row.ContractCreditId.Int64
	if contractId == 0 {
		utils.Fail(c, "该记录缺少链上学分ID，无法审核")
		return
	}

	txHash, err := utils.AuditCredit(uint64(contractId), true)
	if err != nil {
		utils.Fail(c, "链上审核失败: "+err.Error())
		return
	}

	userId, _ := c.Get("userId")
	user, _ := model.GetUserById(userId.(uint64))
	auditAdmin := ""
	if user != nil && user.Address.Valid {
		auditAdmin = user.Address.String
	}
	if err := model.UpdateCreditStatus(req.CreditId, "approved", auditAdmin); err != nil {
		utils.Fail(c, "更新状态失败: "+err.Error())
		return
	}
	utils.Success(c, gin.H{"tx_hash": txHash}, "审核通过")
}

// CreditList 学分列表（按角色：学生看自己、教师看自己录入、管理员看全部）
func CreditList(c *gin.Context) {
	userId, _ := c.Get("userId")
	role, _ := c.Get("role")
	user, err := model.GetUserById(userId.(uint64))
	if err != nil || user == nil {
		utils.Fail(c, "用户不存在")
		return
	}

	var list []model.CreditRow
	switch role.(string) {
	case "student":
		if !user.Address.Valid || user.Address.String == "" {
			utils.Success(c, []interface{}{}, "暂无学分")
			return
		}
		list, err = model.GetCreditsByStudentAddress(user.Address.String)
	case "teacher":
		if !user.Address.Valid || user.Address.String == "" {
			utils.Success(c, []interface{}{}, "暂无录入记录")
			return
		}
		list, err = model.GetCreditsByTeacherAddress(user.Address.String)
	case "admin":
		list, err = model.GetAllCredits()
	default:
		utils.Fail(c, "无权限")
		return
	}
	if err != nil {
		utils.Fail(c, "查询失败: "+err.Error())
		return
	}
	utils.Success(c, list, "查询成功")
}

// CreditPending 管理员：待审核列表
func CreditPending(c *gin.Context) {
	list, err := model.GetPendingCredits()
	if err != nil {
		utils.Fail(c, "查询失败: "+err.Error())
		return
	}
	utils.Success(c, list, "查询成功")
}

// CreditSync 链上学分同步到本地（可根据 student_address 拉取并更新状态）
func CreditSync(c *gin.Context) {
	// 简单实现：把本地 pending 的根据 contract_credit_id 调合约 getCreditById 更新 is_approved 到 status
	list, err := model.GetPendingCredits()
	if err != nil {
		utils.Fail(c, "查询待同步记录失败: "+err.Error())
		return
	}
	// 若没有 GetCreditByIdFromChain，可先只返回成功，后续再实现
	_ = list
	utils.Success(c, gin.H{"synced": 0}, "同步完成（当前仅落库记录，链上状态需审核后更新）")
}

// CreditRejectReq 驳回请求（可选）
type CreditRejectReq struct {
	CreditId int64 `json:"credit_id" binding:"required"`
}

func CreditReject(c *gin.Context) {
	var req CreditRejectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, "参数错误: "+err.Error())
		return
	}
	row, err := model.GetCreditById(req.CreditId)
	if err != nil || row == nil {
		utils.Fail(c, "学分记录不存在")
		return
	}
	if row.Status != "pending" {
		utils.Fail(c, "该记录已处理")
		return
	}
	userId, _ := c.Get("userId")
	user, _ := model.GetUserById(userId.(uint64))
	auditAdmin := ""
	if user != nil && user.Address.Valid {
		auditAdmin = user.Address.String
	}
	if err := model.UpdateCreditStatus(req.CreditId, "rejected", auditAdmin); err != nil {
		utils.Fail(c, "更新失败: "+err.Error())
		return
	}
	utils.Success(c, nil, "已驳回")
}

