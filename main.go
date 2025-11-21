package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/clutchtechnology/hisense-vmi-dataserver/src/databases"
	"github.com/clutchtechnology/hisense-vmi-dataserver/src/models"
	"github.com/clutchtechnology/hisense-vmi-dataserver/src/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if strings.ToLower(os.Getenv("PRODUCTION")) == "true" {
		time.Sleep(10 * time.Second)
	} else {
		if err := godotenv.Load(); err != nil {
			fmt.Println("Error loading .env file:", err)
			return
		}
	}

	// Init DB
	DB_CONN = databases.InitDB(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	models.Migrate(DB_CONN)
	checkAdmin(DB_CONN)

	// Log - 只输出到控制台，不写入文件
	// f, err := os.Create(os.Getenv("LOG_FILE"))
	// if err != nil {
	// 	fmt.Println("LOG_FILE error", err)
	// 	return
	// }
	// defer f.Close()

	// f, err = os.OpenFile(os.Getenv("LOG_FILE"), os.O_WRONLY|os.O_APPEND, os.ModeAppend)
	// if err != nil {
	// 	fmt.Println("err", err)
	// }
	// gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	gin.DefaultWriter = os.Stdout // 只输出到控制台

	// Init godi
	InitGodi()

	// Init gin
	r := gin.Default()
	routes.RegisterRoute(r, SERVICE_CONTAINER)
	r.Run(fmt.Sprintf("%s:%s", os.Getenv("HTTP_HOST"), os.Getenv("HTTP_PORT")))
}
