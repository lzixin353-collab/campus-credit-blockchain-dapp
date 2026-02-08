package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户表
type User struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Address   string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"address"` // 以太坊钱包地址（唯一）
	Role      string         `gorm:"type:varchar(32);default:'student'" json:"role"`       // 角色：admin/teacher/student
	Username  string         `gorm:"type:varchar(64);not null" json:"username"`            // 用户名
	Phone     string         `gorm:"type:varchar(20);nullable" json:"phone"`               // 手机号
	Email     string         `gorm:"type:varchar(64);nullable" json:"email"`               // 邮箱
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`                     // 创建时间
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`                     // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                                       // 软删除
}

// TableName 指定表名
func (u *User) TableName() string {
	return "users"
}
