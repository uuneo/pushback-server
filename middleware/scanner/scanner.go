package scanner

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pushbackServer/middleware/ipban"
	"strings"
)

// 可疑的路径和参数关键词
var suspiciousPatterns = []string{
	// PHP相关
	".php", "wordpress", "wp-content", "wp-includes",
	"wp-admin", "wp-login", "wp-config", "wp-settings",
	"wp-cron", "wp-json", "wp-uploads", "wp-plugins",
	"wp-themes", "wp-signup", "wp-register", "wp-activate",
	"wp-comments", "wp-feed", "wp-rss", "wp-trackback",
	"wp-pingback", "wp-mail", "wp-cron", "wp-links-opml",
	"wp-load", "wp-blog-header", "wp-settings", "wp-cron",
	"wp-comments-post", "wp-login.php", "wp-register.php",
	"wp-signup.php", "wp-activate.php", "wp-comments.php",
	"wp-feed.php", "wp-rss.php", "wp-trackback.php",
	"wp-pingback.php", "wp-mail.php", "wp-cron.php",
	"wp-links-opml.php", "wp-load.php", "wp-blog-header.php",
	"wp-settings.php", "wp-comments-post.php",
}

// CheckIpMiddleware 检测并封禁可疑的扫描请求
func CheckIpMiddleware() gin.HandlerFunc {
	banManager := ipban.GetManager()

	return func(c *gin.Context) {
		// 获取客户端IP
		clientIP := c.ClientIP()

		// 检查IP是否已被封禁
		if banManager.IsBanned(clientIP) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		// 检查请求路径
		path := strings.ToLower(c.Request.URL.Path)
		query := strings.ToLower(c.Request.URL.RawQuery)

		// 检查是否包含可疑特征
		for _, pattern := range suspiciousPatterns {
			if strings.Contains(path, pattern) || strings.Contains(query, pattern) {
				// 封禁IP
				banManager.BanIP(clientIP)
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		}

		// 检查User-Agent
		userAgent := strings.ToLower(c.GetHeader("User-Agent"))
		for _, pattern := range suspiciousPatterns {
			if strings.Contains(userAgent, pattern) {
				banManager.BanIP(clientIP)
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		}

		c.Next()
	}
}
