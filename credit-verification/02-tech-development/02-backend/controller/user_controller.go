// controller/user_controller.go
package controller

import (
	"campus-credit-backend/model"
	"campus-credit-backend/utils"
	"database/sql"

	"github.com/gin-gonic/gin"
)

// UserLogin 用户登录接口
// @Summary 用户登录
// @Description 用户名密码登录，返回JWT Token
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body struct{Username string `json:"username" binding:"required"`;Password string `json:"password" binding:"required"`} true "登录信息"
// @Success 200 {object} utils.Response{data=struct{Token string;User model.User}}
// @Failure 400 {object} utils.Response
// @Router /api/user/login [post]
func UserLogin(c *gin.Context) {
	// 解析请求体
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, "参数错误: "+err.Error())
		return
	}

	// 查询用户
	user, err := model.GetUserByUsername(req.Username)
	if err != nil {
		utils.Fail(c, "查询用户失败: "+err.Error())
		return
	}
	if user == nil || !utils.CheckPassword(req.Password, user.Password) {
		utils.Fail(c, "用户名或密码错误")
		return
	}

	// 生成Token
	token, err := utils.GenerateToken(user.Id, user.Username, user.Role)
	if err != nil {
		utils.Fail(c, "生成Token失败: "+err.Error())
		return
	}

	// 返回结果
	utils.Success(c, gin.H{
		"token": token,
		"user":  user,
	}, "登录成功")
}

// UserInfo 获取当前用户信息
// @Summary 获取用户信息
// @Description 需登录，返回当前用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Header 200 {string} Authorization "Bearer Token"
// @Success 200 {object} utils.Response{data=model.User}
// @Failure 401 {object} utils.Response
// @Router /api/user/info [get]
func UserInfo(c *gin.Context) {
	// 从上下文获取用户ID（中间件注入）
	userId, _ := c.Get("userId")
	user, err := model.GetUserById(userId.(uint64))
	if err != nil {
		utils.Fail(c, "查询用户信息失败: "+err.Error())
		return
	}
	if user == nil {
		utils.Fail(c, "用户不存在")
		return
	}
	utils.Success(c, user, "查询成功")
}

// UserUpdate 更新用户信息
// @Summary 更新用户信息
// @Description 需登录，更新地址/角色（仅管理员可改角色）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Header 200 {string} Authorization "Bearer Token"
// @Param data body struct{Address string `json:"address"`;Role string `json:"role"`} true "更新信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/user/update [post]
func UserUpdate(c *gin.Context) {
	var req struct {
		Address string `json:"address"`
		Role    string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, "参数错误: "+err.Error())
		return
	}

	// 获取当前用户
	userId, _ := c.Get("userId")
	user, err := model.GetUserById(userId.(uint64))
	if err != nil || user == nil {
		utils.Fail(c, "用户不存在")
		return
	}

	// 更新地址
	if req.Address != "" {
		user.Address = sql.NullString{String: req.Address, Valid: true}
	}
	// 仅管理员可更新角色
	currentRole, _ := c.Get("role")
	if currentRole == "admin" && req.Role != "" {
		user.Role = req.Role
	}

	// 执行更新
	if err := user.UpdateUser(); err != nil {
		utils.Fail(c, "更新失败: "+err.Error())
		return
	}
	utils.Success(c, nil, "更新成功")
}

// UserRegister 用户注册（测试用）
// @Summary 用户注册
// @Description 新增用户（测试环境用）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body struct{Username string `json:"username" binding:"required"`;Password string `json:"password" binding:"required"`;Address string `json:"address"`;Role string `json:"role"`} true "注册信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/user/register [post]
func UserRegister(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Address  string `json:"address"`
		Role     string `json:"role" default:"student"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, "参数错误: "+err.Error())
		return
	}

	// 验证角色合法性
	allowedRoles := map[string]bool{"student": true, "teacher": true, "admin": true}
	if !allowedRoles[req.Role] {
		utils.Fail(c, "角色只能是student/teacher/admin")
		return
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Password: req.Password,
		Address:  sql.NullString{String: req.Address, Valid: req.Address != ""},
		Role:     req.Role,
	}
	if err := user.CreateUser(); err != nil {
		utils.Fail(c, "注册失败: "+err.Error())
		return
	}
	utils.Success(c, nil, "注册成功")
}
