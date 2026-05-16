package router

import (
	"fmt"
	"net/http"
	"runtime"

	"eao/internal/config"
	"eao/internal/controller"
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

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "404",
		})
	})
	return r
}
