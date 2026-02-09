// utils/mysql.go
package utils

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DB 全局数据库连接
var DB *sql.DB

// InitMySQL 初始化MySQL连接
func InitMySQL() {
	dsn := GlobalConfig.MySQL.DSN
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 连接池配置
	db.SetMaxOpenConns(GlobalConfig.MySQL.MaxOpenConns)
	db.SetMaxIdleConns(GlobalConfig.MySQL.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(GlobalConfig.MySQL.ConnMaxLifetime) * time.Second)

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatalf("数据库Ping失败: %v", err)
	}
	DB = db
	log.Println("MySQL初始化完成")
}
