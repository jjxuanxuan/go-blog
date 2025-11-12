package routes

import (
	"go-blog/internal/handler"
	"go-blog/internal/middleware"
	"go-blog/internal/model"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	// 自定义中间件链：Recovery -> Logger -> CORS
	router := gin.New()
	router.Use(middleware.Recovery(), middleware.LoggerMiddleware(), middleware.CORS())

	// 初始化数据库（确保 AutoMigrate 已执行）
	model.InitDB()

	uh := handler.NewUserHandler(model.DB)
    ph := handler.NewPostHandler(model.DB)
    ah := handler.NewAuthHandler(model.DB)
    ch := handler.NewCommentHandler(model.DB)

	// 分组：/api/auth
	apiAuth := router.Group("/api/auth")
	{
		apiAuth.POST("/register", ah.Register)
		apiAuth.POST("/login", ah.Login)
		apiAuth.POST("/refresh", ah.Refresh)
	}

	// 分组：/api（鉴权）
    api := router.Group("/api")
    api.Use(middleware.AuthMiddleware(), middleware.RequireUser())
    {
        api.GET("/me", uh.MeHandler)

        api.POST("/posts", ph.CreatePost)
        api.GET("/posts", ph.GetAllPosts)
        api.GET("/posts/:id", ph.GetPostsById)
        api.PUT("/posts/:id", ph.UpdatePost)
        api.DELETE("/posts/:id", ph.DeletePost)

        api.POST("/comments", ch.CreateComment)
        api.DELETE("/comments/:id", ch.DeleteComment)
        api.GET("/posts/:post_id/comments", ch.ListCommentsByPost)
    }

	// 分组：/api/admin（鉴权+RBAC）
	admin := router.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.RequireUser(), middleware.RequireRole("admin"))
	{
		// 预留管理接口
		// 例如：admin.GET("/users", listUsers)
	}

	return router
}
