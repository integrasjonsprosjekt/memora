package router

import (
	"log/slog"
	"memora/internal/config"
	"memora/internal/handlers/docs"
	"memora/internal/handlers/status"
	"memora/internal/middleware"
	"memora/internal/services"

	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	router := gin.New()

	logger := slog.Default()

	router.Use(gin.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.Logging(logger))

	return router
}

func Route(c *gin.Engine, services *services.Services) {
	v1 := c.Group(config.BasePath)
	{
		v1.GET("/status", status.GetStatus)
		v1.GET("/docs", docs.GetDocs)
	}
}
