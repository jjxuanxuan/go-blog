package routes

import (
	"github.com/gin-gonic/gin"
	"go-blog/internal/handler"
	"go-blog/internal/middleware"
	"go-blog/internal/model"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(middleware.LoggerMiddleware())

	// 初始化数据库（确保 AutoMigrate 已执行）
	model.InitDB()

	uh := handler.NewUserHandler(model.DB)

	router.GET("/ping", handler.PingHandler)

	router.POST("/login", uh.LoginHandler)

	router.POST("/register", uh.RegisterHandler)

	return router
}
