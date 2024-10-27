package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pushbackServer/config"
	"pushbackServer/controller"
	"pushbackServer/database"
	"runtime"
)

var (
	version   string
	buildDate string
	commitID  string
)

func Admin() gin.HandlerFunc {
	localUser := config.LocalConfig.System.User
	localPassword := config.LocalConfig.System.Password

	return func(c *gin.Context) {
		// 配置了账号密码，进行身份校验
		if localUser != "" && localPassword != "" {
			// 优先使用 Basic Auth
			user, pass, hasAuth := c.Request.BasicAuth()
			if !hasAuth {
				// 如果没有 Basic Auth，则尝试从查询参数中获取
				user = c.Query("user")
				pass = c.Query("password")

				if user == "" {
					user = c.Query("u")
				}
				if pass == "" {
					pass = c.Query("p")
				}
			}

			if user == localUser && pass == localPassword {
				c.Set("admin", config.LocalConfig.Apple.AdminId)
			}
		} else {
			// 没有配置账号密码，记录 header
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" && len(authHeader) > 1 {
				c.Set("admin", authHeader)
			}
		}

		c.Next()
	}
}

func main() {

	gin.SetMode(config.LocalConfig.System.Mode)
	router := gin.Default()
	router.Use(Admin())
	// 二维码
	router.GET("/", controller.HomeController)
	router.GET("/info", GetInfo)
	// App内部使用
	router.GET("/ping", controller.Ping)
	router.GET("/p", controller.Ping)
	router.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, "ok") })
	router.GET("/h", func(c *gin.Context) { c.JSON(http.StatusOK, "ok") })
	// 注册
	router.POST("/register", controller.RegisterController)
	router.GET("/register/:deviceKey", controller.RegisterController)

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
	addr := config.LocalConfig.System.Host + ":" + config.LocalConfig.System.Post
	if err := router.Run(addr); err != nil {
		panic(err)
	}

}

func GetInfo(c *gin.Context) {
	admin, ok := c.Get("admin")
	if ok && admin.(string) == config.LocalConfig.Apple.AdminId {
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
