package scanner

import (
	"github.com/gin-gonic/gin"
	"pushbackServer/config"
)

func Verification() gin.HandlerFunc {
	localUser := config.LocalConfig.System.User
	localPassword := config.LocalConfig.System.Password

	return func(c *gin.Context) {

		// 先查看是否是管理员身份
		authHeader := c.GetHeader("Authorization")
		if authHeader == config.LocalConfig.Apple.AdminId && authHeader != "" {
			c.Set("admin", true)
			return
		}
		// 配置了账号密码，进行身份校验
		if localUser != "" && localPassword != "" {
			// 优先使用 Basic Auth
			user, pass, hasAuth := c.Request.BasicAuth()
			if !hasAuth {
				// 如果没有 Basic Auth，则尝试从查询参数中获取
				user = c.Query("user")
				pass = c.Query("password")

				if c.Request.Method == "POST" {
					if user == "" {
						user = c.PostForm("username")
					}
					if pass == "" {
						pass = c.PostForm("password")
					}
				} else {
					if user == "" {
						user = c.Query("username")
					}

					if user == "" {
						user = c.Query("u")
					}
					if pass == "" {
						pass = c.Query("p")
					}
				}

			}

			if user == localUser && pass == localPassword {
				c.Set("admin", true)
				return
			}

		}

		// 如果没有身份验证信息
		c.Set("admin", false)
		c.Next()
	}
}
