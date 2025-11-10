package routes

import (
    "github.com/gin-gonic/gin"
    "go-blog/internal/handler"
    "go-blog/internal/middleware"
    "go-blog/internal/model"
)

// SetupRouter 构建 Gin 引擎与路由分组。
// 全局中间件顺序：Recovery(兜底异常) -> Logger(访问日志) -> CORS(跨域)
func SetupRouter() *gin.Engine {
    // 自定义中间件链：Recovery -> Logger -> CORS
    router := gin.New()
    router.Use(middleware.Recovery(), middleware.LoggerMiddleware(), middleware.CORS())

    // 初始化数据库（确保 AutoMigrate 已执行）
    model.InitDB()

    uh := handler.NewUserHandler(model.DB)
    ph := handler.NewPostHandler(model.DB)
    ah := handler.NewAuthHandler(model.DB)

    // 旧接口已删除，统一使用 /api/auth/*
    // 新版分组：/api/auth
    apiAuth := router.Group("/api/auth")
    {
        apiAuth.POST("/register", ah.Register)
        apiAuth.POST("/login", ah.Login)
        apiAuth.POST("/refresh", ah.Refresh)
    }

    // 新版分组：/api（鉴权）—— 需要携带 Access Token
    api := router.Group("/api")
    api.Use(middleware.AuthMiddleware(), middleware.RequireUser())
    {
        api.GET("/me", uh.MeHandler)

        api.POST("/posts", ph.CreatePost)
        api.GET("/posts", ph.GetAllPosts)
        api.GET("/posts/:id", ph.GetPostsById)
        api.PUT("/posts/:id", ph.UpdatePost)
        api.DELETE("/posts/:id", ph.DeletePost)
    }

    // 新版分组：/api/admin（鉴权+RBAC）—— 仅 admin 角色可访问
    admin := router.Group("/api/admin")
    admin.Use(middleware.AuthMiddleware(), middleware.RequireUser(), middleware.RequireRole("admin"))
    {
        // 预留管理接口
        // 例如：admin.GET("/users", listUsers)
    }

    return router
}
