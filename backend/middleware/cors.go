package middleware

import (
	"smtp-mail/backend/config"
	"strings"

	"github.com/gin-gonic/gin"
)

// isLocalhost 检查来源是否为localhost或127.0.0.1
func isLocalhost(origin string) bool {
	return strings.HasPrefix(origin, "http://localhost") ||
		strings.HasPrefix(origin, "http://127.0.0.1")
}

// matchOriginWithWildcard 使用通配符模式匹配来源
func matchOriginWithWildcard(origin string, pattern string) bool {
	// 精确匹配
	if origin == pattern {
		return true
	}

	// 处理通配符模式，如 http://localhost:* 或 http://127.0.0.1:*
	if strings.HasSuffix(pattern, ":*") {
		prefix := strings.TrimSuffix(pattern, ":*")
		if strings.HasPrefix(origin, prefix) {
			return true
		}
	}

	// 处理完整的通配符
	if pattern == "*" {
		return true
	}

	return false
}

// CORS 跨域资源共享中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := config.GetConfig()

		// 设置允许的源
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			allowed := false

			// 检查源是否在允许列表中
			for _, allowedOrigin := range cfg.Security.CORSOrigins {
				if matchOriginWithWildcard(origin, allowedOrigin) {
					allowed = true
					break
				}
			}

			// 如果不在允许列表中，自动允许localhost来源
			if !allowed && isLocalhost(origin) {
				allowed = true
			}

			if allowed {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			}
		}

		// 设置允许的方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// 设置允许的头
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// 设置允许凭证
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// 设置预检请求缓存时间
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		// 处理OPTIONS预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}