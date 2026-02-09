// model/user.go
package model

import (
	"database/sql"
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
