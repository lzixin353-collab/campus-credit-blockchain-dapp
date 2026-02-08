package model

import (
	"time"

	"gorm.io/gorm"
)

// Credit 学分表
type Credit struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint64         `gorm:"index;not null" json:"user_id"`                    // 关联用户ID（仅索引，暂不创建外键）
	CourseName  string         `gorm:"type:varchar(128);not null" json:"course_name"`    // 课程名
	CreditScore float64        `gorm:"type:decimal(5,2);not null" json:"credit_score"`   // 学分值
	Status      string         `gorm:"type:varchar(32);default:'pending'" json:"status"` // 状态：pending/approved/rejected
	Remark      string         `gorm:"type:varchar(256);nullable" json:"remark"`         // 备注
	ContractTx  string         `gorm:"type:varchar(66);nullable" json:"contract_tx"`     // 上链交易哈希
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`                 // 创建时间
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`                 // 更新时间
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`                                   // 软删除

	// 关键修改：注释/删除关联字段的自动外键创建（新手阶段暂不用外键，避免类型问题）
	// User User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}

// TableName 指定表名
func (c *Credit) TableName() string {
	return "credits"
}
