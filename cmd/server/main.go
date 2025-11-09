package main

import (
	"go-blog/internal/routes"
	"log"
)

func main() {

	router := routes.SetupRouter()

	// 启动服务
	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
