package utils

import (
	"campus-credit-backend/model"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitMySQL 初始化MySQL连接
func InitMySQL() {
	dsn := GlobalConfig.MySQL.Dsn
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 显示SQL日志
		// 新增：禁用自动创建外键约束（解决类型不匹配问题）
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("连接MySQL失败: %v", err)
	}

	// 自动迁移表（创建/更新表结构）
	err = db.AutoMigrate(&model.User{}, &model.Credit{})
	if err != nil {
		log.Fatalf("表迁移失败: %v", err)
	}

	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取DB连接池失败: %v", err)
	}
	sqlDB.SetMaxOpenConns(GlobalConfig.MySQL.MaxOpenConns)
	sqlDB.SetMaxIdleConns(GlobalConfig.MySQL.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(GlobalConfig.MySQL.ConnMaxLifetime) * time.Second)

	DB = db
	log.Println("MySQL连接成功，表迁移完成")
}
