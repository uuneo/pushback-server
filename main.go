package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pushbackServer/config"
	"pushbackServer/controller"
)

func main() {
	gin.SetMode(config.LocalConfig.System.Mode)

	router := gin.Default()
	router.Use(Admin())
	router.GET("/", controller.QRCode)
	router.GET("/info", controller.GetInfo)
	// App内部使用
	router.GET("/ping", controller.Ping)
	router.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, "ok") })
	// 注册
	router.POST("/register", controller.RegisterController)
	router.GET("/register/:device_key", controller.RegisterController).Use(Admin())
	// 推送请求
	router.POST("/push", controller.BaseController)
	router.GET("/:deviceKey/:params1/:params2/:params3", controller.BaseController).Use(Auth())
	router.POST("/:deviceKey/:params1/:params2/:params3", controller.BaseController).Use(Auth())
	router.GET("/:deviceKey/:params1/:params2", controller.BaseController).Use(Auth())
	router.POST("/:deviceKey/:params1/:params2", controller.BaseController).Use(Auth())
	router.GET("/:deviceKey/:params1", controller.BaseController).Use(Auth())
	router.POST("/:deviceKey/:params1", controller.BaseController).Use(Auth())
	router.GET("/:deviceKey/pushback", controller.GetPushToken).Use(Auth())
	router.GET("/:deviceKey", controller.BaseController).Use(Auth())
	router.POST("/:deviceKey", controller.BaseController).Use(Auth())

	addr := config.LocalConfig.System.Host + ":" + config.LocalConfig.System.Post
	if err := router.Run(addr); err != nil {
		panic(err)
	}

}

func Auth() gin.HandlerFunc {

	localUser := config.LocalConfig.System.User
	localPassword := config.LocalConfig.System.Password
	if localUser == "" || localPassword == "" {
		return func(c *gin.Context) { c.Next() }
	} else {
		return gin.BasicAuth(gin.Accounts{localUser: localPassword})
	}
}

func Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求头中的 Authorization 字段
		admin := c.GetHeader("Authorization")
		c.Set("admin", admin)
		// 如果通过验证，继续执行请求处理
		c.Next()
	}
}
