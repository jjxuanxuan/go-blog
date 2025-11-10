package model

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func InitDB() {
	appEnv := getEnv("APP_ENV", "development")

	if appEnv != "production" {
		if err := godotenv.Load(); err == nil {
			fmt.Println("已加载 .env配置文件")
		} else {
			fmt.Println("未找到 .env 文件")
		}
	} else {
		fmt.Println("生产环境：从系统环境变量读取配置")
	}

	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnv("DB_PORT", "3306")
	user := getEnv("DB_USER", "app")
	pass := getEnv("DB_PASS", "123456")
	name := getEnv("DB_NAME", "go_blog")

	//DB_DSN=app:123456@tcp(127.0.0.1:3306)/go_blog?charset=utf8mb4&parseTime=true&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		user, pass, host, port, name)

	//GORM 日志配置：开发打印 SQL，生产只打印错误
	var gormLogger logger.Interface
	if appEnv == "production" {
		gormLogger = logger.Default.LogMode(logger.Error)
	} else {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Fatalf("数据库连接失败 %v", err)
	}
	fmt.Println("数据库连接成功")

	//设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取底层 sql.DB 失败: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db

	if err := DB.AutoMigrate(&User{}, &Post{}); err != nil {
		log.Fatalf("auto migrate error: %v", err)
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
