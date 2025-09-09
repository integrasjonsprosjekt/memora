package router

import (
	"memora/internal/config"
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

	v1 := router.Group("/api/v1")
	{
		v1.GET("/status", status.GetStatus)
	}

	return router
}
