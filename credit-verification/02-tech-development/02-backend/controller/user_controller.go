// controller/user_controller.go
package controller

import (
	"campus-credit-backend/model"
	"campus-credit-backend/utils"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

// UserLogin 用户登录接口（支持用户名密码 + 钱包地址两种方式）
// @Summary 用户登录
// @Description 传 address 为钱包登录（从链上取角色）；传 username+password 为账号密码登录
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body struct{Address string `json:"address"`;Username string `json:"username"`;Password string `json:"password"`} true "登录信息"
// @Success 200 {object} utils.Response{data=struct{Token string;User model.User}}
// @Failure 400 {object} utils.Response
// @Router /api/user/login [post]
func UserLogin(c *gin.Context) {
	// 最外层：任何 panic 都返回 200+JSON，避免 500 空 body
	defer func() {
		if r := recover(); r != nil && !c.Writer.Written() {
			log.Printf("[UserLogin] panic: %v", r)
			utils.Fail(c, "登录异常: "+panicToStr(r))
		}
	}()
	var req struct {
		Address  string `json:"address"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, "参数错误: "+err.Error())
		return
	}

	// 钱包登录：仅传 address
	if req.Address != "" {
		userLoginByAddress(c, req.Address)
		return
	}

	// 账号密码登录
	if req.Username == "" || req.Password == "" {
		utils.Fail(c, "请提供 username 与 password，或提供 address 进行钱包登录")
		return
	}
	user, err := model.GetUserByUsername(req.Username)
	if err != nil {
		utils.Fail(c, "查询用户失败: "+err.Error())
		return
	}
	if user == nil || !utils.CheckPassword(req.Password, user.Password) {
		utils.Fail(c, "用户名或密码错误")
		return
	}
	token, err := utils.GenerateToken(user.Id, user.Username, user.Role)
	if err != nil {
		utils.Fail(c, "生成Token失败: "+err.Error())
		return
	}
	utils.Success(c, gin.H{
		"token": token,
		"user":  user,
	}, "登录成功")
}

func panicToStr(r interface{}) string {
	switch x := r.(type) {
	case string:
		return x
	case error:
		return x.Error()
	default:
		return fmt.Sprintf("%v", r)
	}
}

// normalAddress 钱包地址规范化，便于与 DB 匹配（去空格、小写）
func normalAddress(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// userLoginByAddress 钱包地址登录：从链上取角色，无则自动建用户
func userLoginByAddress(c *gin.Context, address string) {
	defer func() {
		if r := recover(); r != nil && !c.Writer.Written() {
			log.Printf("[userLoginByAddress] panic: %v", r)
			utils.Fail(c, "登录处理异常: "+panicToStr(r))
		}
	}()
	addr := normalAddress(address)
	if addr == "" {
		utils.Fail(c, "钱包地址无效")
		return
	}
	user, err := model.GetUserByAddress(addr)
	if err != nil {
		utils.Fail(c, "查询用户失败: "+err.Error())
		return
	}
	if user != nil {
		// 数据库已有该地址：直接使用 DB 角色，不按链上覆盖（避免 Postman 注册的 admin 被链上未分配而变成 student）
		token, err := utils.GenerateToken(user.Id, user.Username, user.Role)
		if err != nil {
			utils.Fail(c, "生成Token失败: "+err.Error())
			return
		}
		utils.Success(c, gin.H{"token": token, "user": user}, "登录成功")
		return
	}
	// 新钱包用户：从链上取角色再建用户，链上无则默认 student
	role, err := utils.GetRoleFromChain(addr)
	if err != nil {
		utils.Fail(c, "获取链上角色失败: "+err.Error())
		return
	}
	user, err = model.CreateWalletUser(addr, role)
	if err != nil {
		utils.Fail(c, "创建用户失败: "+err.Error())
		return
	}
	if user == nil {
		utils.Fail(c, "创建用户失败，请重试")
		return
	}
	token, err := utils.GenerateToken(user.Id, user.Username, user.Role)
	if err != nil {
		utils.Fail(c, "生成Token失败: "+err.Error())
		return
	}
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

	// 更新地址（规范化后入库，便于钱包登录时匹配）
	if req.Address != "" {
		user.Address = sql.NullString{String: normalAddress(req.Address), Valid: true}
	}
	// 仅管理员可更新角色
	currentRole, _ := c.Get("role")
	if r, ok := currentRole.(string); ok && strings.ToLower(r) == "admin" && req.Role != "" {
		user.Role = strings.ToLower(strings.TrimSpace(req.Role))
	}

	// 执行更新
	if err := user.UpdateUser(); err != nil {
		utils.Fail(c, "更新失败: "+err.Error())
		return
	}
	utils.Success(c, nil, "更新成功")
}

// BindAddress 绑定钱包地址到当前账号（先释放被 wallet_ 占用的同地址，再绑定）
// @Summary 绑定钱包地址
// @Description 登录后调用，将钱包地址绑到当前账号；若该地址已被自动创建的 wallet_ 占用则先删除再绑定
// @Tags 用户管理
// @Accept json
// @Produce json
// @Header 200 {string} Authorization "Bearer Token"
// @Param data body struct{Address string `json:"address" binding:"required"`} true "钱包地址"
// @Success 200 {object} utils.Response
// @Router /api/user/bind-address [post]
func BindAddress(c *gin.Context) {
	var req struct {
		Address string `json:"address" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, "请提供 address")
		return
	}
	addr := normalAddress(req.Address)
	if addr == "" {
		utils.Fail(c, "钱包地址无效")
		return
	}

	userId, _ := c.Get("userId")
	user, err := model.GetUserById(userId.(uint64))
	if err != nil || user == nil {
		utils.Fail(c, "用户不存在")
		return
	}

	// 已绑定该地址
	if user.Address.Valid && normalAddress(user.Address.String) == addr {
		utils.Success(c, nil, "已绑定该地址")
		return
	}

	// 若该地址被 wallet_ 用户占用，先删除以便当前用户可绑定（表有 UNIQUE(address)）
	_, _ = model.FreeAddressFromWalletUser(addr)

	// 检查是否被其他非 wallet_ 账号占用
	other, _ := model.GetUserByAddress(addr)
	if other != nil && other.Id != user.Id {
		utils.Fail(c, "该地址已被其他账号绑定")
		return
	}

	user.Address = sql.NullString{String: addr, Valid: true}
	if err := user.UpdateAddress(addr); err != nil {
		utils.Fail(c, "绑定失败: "+err.Error())
		return
	}
	utils.Success(c, nil, "绑定成功，可使用该钱包登录并保持当前角色")
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

	// 角色转小写并验证
	roleLower := strings.ToLower(strings.TrimSpace(req.Role))
	if roleLower == "" {
		roleLower = "student"
	}
	allowedRoles := map[string]bool{"student": true, "teacher": true, "admin": true}
	if !allowedRoles[roleLower] {
		utils.Fail(c, "角色只能是 student/teacher/admin")
		return
	}

	// 创建用户（角色、地址统一规范化入库）
	addrVal := ""
	if req.Address != "" {
		addrVal = normalAddress(req.Address)
	}
	user := &model.User{
		Username: req.Username,
		Password: req.Password,
		Address:  sql.NullString{String: addrVal, Valid: addrVal != ""},
		Role:     roleLower,
	}
	if err := user.CreateUser(); err != nil {
		utils.Fail(c, "注册失败: "+err.Error())
		return
	}
	utils.Success(c, nil, "注册成功")
}
