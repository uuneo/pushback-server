package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"pushbackServer/config"
	"pushbackServer/controller"
	"pushbackServer/database"
	"pushbackServer/middleware/cache"
	"pushbackServer/middleware/scanner"
	"runtime"
)

var (
	version   string
	buildDate string
	commitID  string
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	if config.LocalConfig.System.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	if config.LocalConfig.System.Debug {
		router = gin.Default()
		gin.ForceConsoleColor()
	} else {
		router.Use(gin.Recovery())
	}

	// 添加禁用缓存中间件（在最前面）
	router.Use(cache.NoCache())

	// 添加扫描器检测中间件
	router.Use(scanner.CheckIpMiddleware())
	router.Use(scanner.Verification())

	tmpl := template.Must(template.New("").ParseFS(staticFiles, "static/*.html"))
	router.SetHTMLTemplate(tmpl)

	// 二维码
	router.GET("/", controller.HomeController)
	router.GET("/info", GetInfo)

	router.GET("/upload", controller.UploadController)
	router.GET("/u", controller.UploadController)
	router.GET("/image/:filename", controller.GetImage)
	router.GET("/img/:filename", controller.GetImage)

	// App内部使用
	router.GET("/ping", controller.Ping)
	router.GET("/p", controller.Ping)
	router.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, "ok") })
	router.GET("/h", func(c *gin.Context) { c.JSON(http.StatusOK, "ok") })
	// 注册
	router.POST("/register", controller.RegisterController)
	router.GET("/register/:deviceKey", controller.RegisterController)

	router.POST("/upload", controller.UploadController)
	router.POST("/u", controller.UploadController)

	// 推送请求
	router.POST("/push", controller.BaseController)
	router.POST("/p", controller.BaseController)

	// title subtitle body
	router.GET("/:deviceKey/:params1/:params2/:params3", controller.BaseController)
	router.POST("/:deviceKey/:params1/:params2/:params3", controller.BaseController)
	// title body
	router.GET("/:deviceKey/:params1/:params2", controller.BaseController)
	router.POST("/:deviceKey/:params1/:params2", controller.BaseController)
	// body
	router.GET("/:deviceKey/:params1", controller.BaseController)
	router.POST("/:deviceKey/:params1", controller.BaseController)
	// 获取设备Token
	router.GET("/:deviceKey/token", controller.GetPushToken)
	router.GET("/:deviceKey/t", controller.GetPushToken)
	// 静默推送
	router.GET("/:deviceKey/update", controller.SilentPushController)
	router.GET("/:deviceKey/u", controller.SilentPushController)
	// 参数化的推送
	router.GET("/:deviceKey", controller.BaseController)
	router.POST("/:deviceKey", controller.BaseController)
	log.Println("listening on:", config.LocalConfig.System.Address)
	if err := router.Run(config.LocalConfig.System.Address); err != nil {
		panic(err)
	}

}

func GetInfo(c *gin.Context) {
	admin, ok := c.Get("admin")
	if ok && admin.(bool) {
		devices, _ := database.DB.CountAll()
		c.JSON(http.StatusOK, map[string]interface{}{
			"version": version,
			"build":   buildDate,
			"commit":  commitID,
			"devices": devices,
			"arch":    runtime.GOOS + "/" + runtime.GOARCH,
			"cpu":     runtime.NumCPU(),
		})
	} else {
		c.JSON(http.StatusOK, map[string]interface{}{
			"version": version,
			"build":   buildDate,
			"commit":  commitID,
			"arch":    runtime.GOOS + "/" + runtime.GOARCH,
		})
	}

}
