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

	router.Use(gin.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.Logging())

	return router
}

func Route(c *gin.Engine) {
	v1 := c.Group(config.BasePath)
	{
		v1.GET("/status", status.GetStatus)
		v1.GET("/docs", docs.GetDocs)
	}
}
