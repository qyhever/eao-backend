package router

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"eao/internal/config"
	"eao/internal/controller"
	"eao/internal/middleware"
	"eao/internal/model"
	dbpkg "eao/internal/pkg/db"
	"eao/internal/pkg/password"
	"eao/internal/repository"
	"eao/internal/repository/persistence"
	"eao/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	isProd := config.IsProduction()
	// Gin 开启生产模式(默认是debug模式，会输出大量调试日志)
	if isProd {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// 静态文件服务
	r.Static("/public", "./public")

	fmt.Printf("Go Version %v\n", runtime.Version())

	cfg := config.GetConfig()
	db, err := dbpkg.OpenMySQLFromConfig(cfg)
	if err != nil {
		panic(fmt.Errorf("初始化 MySQL 失败: %w", err))
	}
	adminRepo := newAdminRepositoryFromConfig(cfg, db)
	adminService := service.NewAdminAccountService(adminRepo)
	adminAccountController := controller.NewAdminAccountController(adminService)
	adminAuthService := service.NewAdminAuthService(adminRepo)
	adminAuthController := controller.NewAdminAuthController(adminAuthService)

	metaController := controller.NewMetaController()
	appRepo := persistence.NewAppRepository()
	appService := service.NewAppService(appRepo)
	appController := controller.NewAppController(appService)

	postRepo := persistence.NewPostRepository()
	postService := service.NewPostService(postRepo)
	postController := controller.NewPostController(postService)

	videoRepo := persistence.NewVideoRepository()
	videoService := service.NewVideoService(videoRepo)
	videoController := controller.NewVideoController(videoService)

	fileRepo := persistence.NewFileRepository(cfg.ThirdParty.FileAPI)
	fileService := service.NewFileService(fileRepo)
	fileController := controller.NewFileController(fileService)

	v1 := r.Group("/api")

	v1.GET("/meta", metaController.GetMeta)

	commonGroup := v1.Group("/common")
	appGroup := v1.Group("/app")
	adminGroup := v1.Group("/admin")

	{
		appGroup.POST("/getHelloInfo", appController.GetHelloInfo)
	}

	postGroup := commonGroup.Group("/post")
	{
		postGroup.GET("", postController.GetPostList)
		postGroup.GET("/:id", postController.GetPostByID)
		postGroup.POST("", postController.CreatePost)
		postGroup.PUT("/:id", postController.UpdatePost)
		postGroup.DELETE("/:id", postController.DeletePost)
	}

	fileGroup := commonGroup.Group("/file")
	{
		fileGroup.GET("/list", fileController.List)
		fileGroup.GET("/listByDir", fileController.ListByDir)
		fileGroup.POST("/upload", fileController.Upload)
	}

	video := v1.Group("/video")
	{
		video.GET("", videoController.GetVideoList)
	}

	adminGroup.POST("/auth/login", adminAuthController.AdminLogin)
	adminGroup.POST("/auth/refresh", adminAuthController.AdminRefreshToken)
	adminProtectedGroup := adminGroup.Group("")
	adminProtectedGroup.Use(middleware.JWTAuthMiddleware())
	{
		adminProtectedGroup.GET("/users/:id", adminAccountController.GetAdmin)
		adminProtectedGroup.PUT("/users/:id", adminAccountController.UpdateAdmin)
		adminProtectedGroup.DELETE("/users/batch", adminAccountController.BatchDeleteAdmins)
		adminProtectedGroup.PUT("/users/:id/status", adminAccountController.ToggleAdminStatus)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "404",
		})
	})
	return r
}

func newAdminRepositoryFromConfig(cfg *config.Config, db *sql.DB) repository.AdminAccountRepository {
	admin, err := buildAdminSeed(cfg)
	if err != nil {
		panic(fmt.Errorf("初始化管理员 seed 失败: %w", err))
	}
	repo := persistence.NewAdminAccountRepository(db)
	if admin == nil {
		return repo
	}

	if err := repo.Upsert(context.Background(), *admin); err != nil {
		panic(fmt.Errorf("写入管理员 seed 失败: %w", err))
	}

	return repo
}

func buildAdminSeed(cfg *config.Config) (*model.Admin, error) {
	if cfg == nil {
		return nil, nil
	}

	adminCfg := cfg.Auth.Admin
	if strings.TrimSpace(adminCfg.Username) == "" || strings.TrimSpace(adminCfg.Password) == "" {
		return nil, nil
	}

	hash, err := password.Hash(adminCfg.Password)
	if err != nil {
		return nil, err
	}

	name := adminCfg.Name
	if strings.TrimSpace(name) == "" {
		name = adminCfg.Username
	}

	return &model.Admin{
		ID:           1,
		Username:     adminCfg.Username,
		PasswordHash: hash,
		Name:         name,
		Status:       "active",
	}, nil
}
