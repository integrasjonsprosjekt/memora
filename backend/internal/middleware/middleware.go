package middleware

import (
	"log/slog"
	"memora/internal/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// Logging logs details about each request and its outcome
// using the provided slog.Logger instance.
func Logging(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Generate a unique request ID and set it in the context.
		reqID := uuid.New().String()
		c.Set("reqID", reqID)

		c.Next()

		duration := time.Since(start)

		// Log based on the configured log level.
		switch config.CurrentLevel {
		case config.LogLevelDebug:
			logger.Debug("Request completed",
				"reqID", reqID,
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"query", c.Request.URL.RawQuery,
				"status", c.Writer.Status(),
				"duration", duration,
				"client_ip", c.ClientIP(),
				"user_agent", c.Request.UserAgent(),
				"response_size", c.Writer.Size(),
				"errors", c.Errors.String(),
			)
		case config.LogLevelInfo:
			logger.Info("Request completed",
				"reqID", reqID,
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"query", c.Request.URL.RawQuery,
				"status", c.Writer.Status(),
				"duration", duration,
			)
		case config.LogLevelError:
			for _, e := range c.Errors {
				logger.Error("Request error", "reqID", reqID, "error", e.Error())
			}
		}
	}
}
