package middleware

import (
	"log"
	"memora/internal/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CORS adds CORS-related headers to the response
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		log.Printf("Got request: %s %s %d %s %v", c.Request.Method, c.Request.URL.Path, status, clientIP, duration)

		if config.CurrentLevel == config.LogLevelDebug {
			log.Println("Request details")
			for name, values := range c.Request.Header {
				for _, value := range values {
					log.Printf("Header: %s=%s", name, value)
				}
			}
		}
	}
}
