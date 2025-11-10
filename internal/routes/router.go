package routes

import (
    "github.com/gin-gonic/gin"
    "go-blog/internal/handler"
    "go-blog/internal/middleware"
    "go-blog/internal/model"
)

func SetupRouter() *gin.Engine {
    // 使用 Gin 自带的 Logger + Recovery
    router := gin.Default()

	// 初始化数据库（确保 AutoMigrate 已执行）
	model.InitDB()

	uh := handler.NewUserHandler(model.DB)
	ph := handler.NewPostHandler(model.DB)

	router.POST("/login", uh.UserLogin)
	router.POST("/register", uh.UserRegister)

	auth := router.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/me", uh.MeHandler)

		auth.POST("/posts", ph.CreatePost)       // 创建文章
		auth.GET("/posts", ph.GetAllPosts)       //获取所有文章列表
		auth.GET("/posts/:id", ph.GetPostsById)  //查看单个文章详情
		auth.PUT("/posts/:id", ph.UpdatePost)    //更新文章内容
		auth.DELETE("/posts/:id", ph.DeletePost) //删除文章
	}
	return router
}
