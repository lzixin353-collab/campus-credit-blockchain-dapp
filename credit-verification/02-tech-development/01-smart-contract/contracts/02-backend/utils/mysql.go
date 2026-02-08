package utils

import (
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
	})
	if err != nil {
		log.Fatalf("连接MySQL失败: %v", err)
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
	log.Println("MySQL连接成功")
}
