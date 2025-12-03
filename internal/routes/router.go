package routes

import (
	"go-blog/internal/handler"
	"go-blog/internal/middleware"
	"go-blog/internal/model"
	"go-blog/internal/repository"
	"go-blog/internal/service"

	"github.com/gin-gonic/gin"
)

const uploadRoot = "storage/uploads"

// SetupRouter 初始化 Gin 路由、中间件及依赖注入。
func SetupRouter() *gin.Engine {
	// 自定义中间件链：Recovery -> Logger -> CORS
	router := gin.New()
	router.Use(middleware.Recovery(), middleware.LoggerMiddleware(), middleware.CORS())

	// 初始化数据库（确保 AutoMigrate 已执行）
	model.InitDB()

	userRepo := repository.NewUserRepository(model.DB)
	postRepo := repository.NewPostRepository(model.DB)
	userSvc := service.NewUserService(userRepo, postRepo)
	postSvc := service.NewPostService(model.DB, postRepo)
	authSvc := service.NewAuthService(userRepo)
	commentRepo := repository.NewCommentRepository(model.DB)
	categoryRepo := repository.NewCategoryRepository(model.DB)
	tagRepo := repository.NewTagRepository(model.DB)
	uploadRepo := repository.NewUploadRepository(uploadRoot)
	commentSvc := service.NewCommentService(commentRepo, postRepo)
	categorySvc := service.NewCategoryService(categoryRepo)
	tagSvc := service.NewTagService(tagRepo)
	uploadSvc := service.NewUploadService(uploadRepo)
	adminSvc := service.NewAdminService(userRepo, postRepo, commentRepo)

	uh := handler.NewUserHandler(userSvc)
	ph := handler.NewPostHandler(postSvc)
	ah := handler.NewAuthHandler(authSvc)
	ch := handler.NewCommentHandler(commentSvc)
	gh := handler.NewCategoryHandler(categorySvc)
	th := handler.NewTagHandler(tagSvc)
	fh := handler.NewUploadHandler(uploadSvc)
	adh := handler.NewAdminHandler(adminSvc)

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
		api.POST("/comments/:id/reply", ch.ReplyComment)
		api.DELETE("/comments/:id", ch.DeleteComment)
		api.GET("/posts/:id/comments", ch.ListCommentsByPost)

		api.GET("/users/:id/posts", uh.ListUserPosts)

		api.GET("/categories", gh.ListCategories)
		api.POST("/categories", gh.CreateCategory)

		api.GET("/tags", th.ListTags)
		api.POST("/tags", th.CreateTag)

		api.POST("/upload", fh.UploadSingle)
		api.POST("/upload/multi", fh.UploadMulti)
	}
	router.Static("/static/uploads", "./"+uploadRoot)

	// 分组：/api/admin（鉴权+RBAC）
	admin := router.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.RequireUser(), middleware.RequireRole("admin"))
	{
		admin.GET("/dashboard", adh.Dashboard)
		admin.GET("/users", adh.ListUsers)
		admin.GET("/posts", adh.ListPosts)
		admin.GET("/comments", adh.ListComments)
	}

	return router
}
