package databases

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(host string, user string, password string, port string, dbName string) *gorm.DB {
	// Connect to DB
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	// 确保使用 UTF-8 字符集（SET NAMES 包含了 client、connection、results 三个变量）
	db.Exec("SET NAMES utf8mb4")

	return db
}
