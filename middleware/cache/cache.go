package cache

import "github.com/gin-gonic/gin"

// NoCache 中间件用于禁用浏览器缓存
func NoCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置禁用缓存的响应头
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Header("X-Accel-Buffering", "no") // 禁用 Nginx 缓冲
		
		// 对于静态文件，添加版本号或时间戳
		if c.Request.URL.Path == "/upload" || c.Request.URL.Path == "/u" {
			c.Header("Last-Modified", "0")
		}
		
		c.Next()
	}
} 