package router

import (
	"memora/internal/config"
	"memora/internal/handlers/docs"
	"memora/internal/handlers/status"
	"memora/internal/middleware"

	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	router := gin.New()

	if config.CurrentLevel == config.LogLevelDebug ||
		config.CurrentLevel == config.LogLevelInfo {
		router.Use(middleware.Logging())
	}

	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	return router
}

func Route(c *gin.Engine) {
	v1 := c.Group("/api/v1")
	{
		v1.GET("/status", status.GetStatus)
	}
	{
		v1.GET("/docs", docs.GetDocs)
	}
}
