// model/user.go
package model

import (
	"database/sql"
	"fmt"
	"time"

	"campus-credit-backend/utils"
)

// User 用户模型（对应users表）
type User struct {
	Id        uint64         `json:"id"`
	Username  string         `json:"username"`
	Password  string         `json:"-"`       // 序列化时隐藏密码
	Address   sql.NullString `json:"address"` // 以太坊地址（允许为空）
	Role      string         `json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// CreateUser 创建用户（注册）
func (u *User) CreateUser() error {
	// 密码加密
	hashPwd, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashPwd

	// 插入数据库
	_, err = utils.DB.Exec(
		"INSERT INTO users (username, password, address, role) VALUES (?, ?, ?, ?)",
		u.Username, u.Password, u.Address, u.Role,
	)
	return err
}

// GetUserByUsername 按用户名查询用户
func GetUserByUsername(username string) (*User, error) {
	var user User
	err := utils.DB.QueryRow(
		"SELECT id, username, password, address, role, created_at, updated_at FROM users WHERE username = ?",
		username,
	).Scan(
		&user.Id, &user.Username, &user.Password, &user.Address,
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil // 无此用户
	}
	return &user, err
}

// GetUserByAddress 按钱包地址查询用户（用于钱包登录，不区分大小写）
// 若同一地址有多条记录，优先返回 admin > teacher > student，避免先创建的 student 覆盖后绑定的 admin
func GetUserByAddress(address string) (*User, error) {
	if address == "" {
		return nil, nil
	}
	var user User
	err := utils.DB.QueryRow(
		`SELECT id, username, password, address, role, created_at, updated_at 
		 FROM users 
		 WHERE LOWER(TRIM(address)) = LOWER(TRIM(?)) 
		 ORDER BY CASE LOWER(TRIM(role)) WHEN 'admin' THEN 1 WHEN 'teacher' THEN 2 ELSE 3 END 
		 LIMIT 1`,
		address,
	).Scan(
		&user.Id, &user.Username, &user.Password, &user.Address,
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

// GetUserById 按ID查询用户
func GetUserById(userId uint64) (*User, error) {
	var user User
	err := utils.DB.QueryRow(
		"SELECT id, username, address, role, created_at, updated_at FROM users WHERE id = ?",
		userId,
	).Scan(
		&user.Id, &user.Username, &user.Address, &user.Role,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

// UpdateUser 更新用户信息
func (u *User) UpdateUser() error {
	_, err := utils.DB.Exec(
		"UPDATE users SET address = ?, role = ?, updated_at = NOW() WHERE id = ?",
		u.Address, u.Role, u.Id,
	)
	return err
}

// UpdateAddress 仅更新当前用户的钱包地址（用于绑定）
func (u *User) UpdateAddress(address string) error {
	_, err := utils.DB.Exec(
		"UPDATE users SET address = ?, updated_at = NOW() WHERE id = ?",
		address, u.Id,
	)
	return err
}

// FreeAddressFromWalletUser 解除被 wallet_ 占用的地址（删除该条以便绑到其他账号）
func FreeAddressFromWalletUser(address string) (int64, error) {
	res, err := utils.DB.Exec(
		"DELETE FROM users WHERE LOWER(TRIM(address)) = LOWER(TRIM(?)) AND username LIKE 'wallet_%'",
		address,
	)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// CreateWalletUser 钱包用户首次登录时创建（仅地址+角色，无密码登录）
func CreateWalletUser(address, role string) (*User, error) {
	username := "wallet_" + address
	if len(username) > 50 {
		username = username[:50]
	}
	// 使用占位密码，避免空字符串导致 bcrypt 异常
	hashPwd, err := utils.HashPassword("wallet-nologin")
	if err != nil {
		return nil, err
	}
	res, err := utils.DB.Exec(
		"INSERT INTO users (username, password, address, role) VALUES (?, ?, ?, ?)",
		username, hashPwd, address, role,
	)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	if id <= 0 {
		return nil, fmt.Errorf("插入用户后未获得有效ID")
	}
	user, err := GetUserById(uint64(id))
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("创建后无法查询到用户")
	}
	return user, nil
}
