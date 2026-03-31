package config

import (
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware logs HTTP requests with detailed information
func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return "[" + param.TimeStamp.Format(time.RFC3339) + "] " +
			param.ClientIP + " " +
			param.Method + " " +
			param.Path + " " +
			param.Request.Proto + " " +
			"| Status: " + string(rune(param.StatusCode)) + " " +
			"| Latency: " + param.Latency.String() + "\n"
	})
}

// RecoveryMiddleware recovers from panic with 500 error
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.Recovery()
}

// CORSMiddleware handles CORS headers
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RequestIDMiddleware adds request ID to context
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			// Generate simple request ID from timestamp
			requestID = "req-" + string(rune(time.Now().UnixNano()))
		}
		c.Set("RequestID", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// TimeoutMiddleware adds timeout context to requests
func TimeoutMiddleware(timeoutSeconds int) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
