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
	// adminRepo := newAdminRepositoryFromConfig(cfg, db)
	newAdminRepositoryFromConfig(cfg, db)

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

	v1 := r.Group("/api")

	v1.GET("/meta", metaController.GetMeta)

	app := v1.Group("/app")
	{
		app.POST("/getHelloInfo", appController.GetHelloInfo)
	}

	post := v1.Group("/post")
	{
		post.GET("", postController.GetPostList)
		post.GET("/:id", postController.GetPostByID)
		post.POST("", postController.CreatePost)
		post.PUT("/:id", postController.UpdatePost)
		post.DELETE("/:id", postController.DeletePost)
	}

	video := v1.Group("/video")
	{
		video.GET("", videoController.GetVideoList)
	}

	// 下面的api是需要登录的
	v1.Use(middleware.JWTAuthMiddleware())

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "404",
		})
	})
	return r
}

func newAdminRepositoryFromConfig(cfg *config.Config, db *sql.DB) repository.AdminRepository {
	admin, err := buildAdminSeed(cfg)
	if err != nil {
		panic(fmt.Errorf("初始化管理员 seed 失败: %w", err))
	}
	repo := persistence.NewAdminRepositoryWithDB(db)
	if admin == nil {
		return repo
	}

	if db != nil {
		if err := repo.Upsert(context.Background(), *admin); err != nil {
			panic(fmt.Errorf("写入管理员 seed 失败: %w", err))
		}
		return repo
	}

	return persistence.NewAdminRepository(*admin)
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
