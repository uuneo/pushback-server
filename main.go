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
	// App内部使用
	router.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, "ok") })
	router.GET("/ping", controller.Ping)
	router.GET("/server", controller.QRCode)
	router.POST("/register", controller.RegisterController)

	router.GET("/info", controller.GetInfo).Use(Auth())
	router.POST("/change", controller.ChangeKeyHandler).Use(Auth())
	// 推送请求
	router.POST("/push", controller.BaseController)
	router.GET("/:deviceKey/:params1/:params2/:params3", controller.BaseController).Use(Auth())
	router.POST("/:deviceKey/:params1/:params2/:params3", controller.BaseController).Use(Auth())
	router.GET("/:deviceKey/:params1/:params2", controller.BaseController).Use(Auth())
	router.POST("/:deviceKey/:params1/:params2", controller.BaseController).Use(Auth())

	router.GET("/:deviceKey/:params1", controller.BaseController).Use(Auth())
	router.POST("/:deviceKey/:params1", controller.BaseController).Use(Auth())
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
		return gin.BasicAuth(gin.Accounts{
			localUser: localPassword,
		})
	}
}
