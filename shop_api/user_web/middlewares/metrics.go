package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RequestMetricsMiddleware 记录请求的性能指标
func RequestMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 计算请求处理时间
		duration := time.Since(startTime)

		// 记录请求信息
		zap.S().Infow("API请求",
			"path", c.Request.URL.Path,
			"method", c.Request.Method,
			"status", c.Writer.Status(),
			"duration", duration.String(),
			"client_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
		)
	}
}

// SecurityHeadersMiddleware 添加安全相关的HTTP头
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 防止网站被嵌入到iframe中，避免点击劫持
		c.Header("X-Frame-Options", "DENY")
		
		// 防止浏览器对内容类型的猜测
		c.Header("X-Content-Type-Options", "nosniff")
		
		// 启用XSS过滤
		c.Header("X-XSS-Protection", "1; mode=block")
		
		// 内容安全策略
		c.Header("Content-Security-Policy", "default-src 'self'")
		
		// 限制引用来源
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		c.Next()
	}
}

// RateLimitMiddleware 简单的限流中间件
func RateLimitMiddleware() gin.HandlerFunc {
	// 这里可以使用Redis或内存实现更复杂的限流逻辑
	// 这只是一个简单示例
	return func(c *gin.Context) {
		// TODO: 实现基于IP或用户ID的限流逻辑
		// 这里可以集成令牌桶或漏桶算法
		
		c.Next()
	}
}

// ErrorHandlerMiddleware 统一错误处理中间件
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录panic错误
				zap.S().Errorw("API服务发生panic",
					"error", err,
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
					"client_ip", c.ClientIP(),
				)
				
				// 返回500错误给客户端
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":  500,
					"error": "服务器内部错误",
				})
				
				// 终止后续中间件
				c.Abort()
			}
		}()
		
		c.Next()
	}
}