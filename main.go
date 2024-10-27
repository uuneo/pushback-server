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
	router.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, "ok") })
	router.GET("/ping", controller.Ping)
	router.GET("/server", controller.QRCode)
	router.POST("/register", controller.RegisterController)
	router.GET("/info", controller.GetInfo).Use(Auth())
	router.POST("/push", controller.BaseController).Use(Auth())
	router.POST("/change", controller.ChangeKeyHandler).Use(Auth())

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
	return func(c *gin.Context) {
		localUser := config.LocalConfig.System.User
		localPassword := config.LocalConfig.System.Password
		if localUser == "" || localPassword == "" {
			c.Next()
			return
		}

		var user, password string

		if c.Request.Method == "GET" {
			user = c.Query("user")
			password = c.Query("password")
		} else if c.Request.Method == "POST" {
			user = c.PostForm("user")
			password = c.PostForm("password")
		}

		if user != localUser || password != localPassword {
			c.JSON(401, gin.H{
				"status":  "failed",
				"message": "I'm a teapot",
			})

			c.Abort()
			return
		}

	}
}
