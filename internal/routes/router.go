package routes

import (
	"github.com/gin-gonic/gin"
	"go-blog/internal/handler"
	"go-blog/internal/middleware"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	
	router.Use(middleware.LoggerMiddleware())

	router.GET("/ping", handler.PingHandler)

	router.POST("/login", handler.LoginHandler)

	return router
}
