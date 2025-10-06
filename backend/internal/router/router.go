package router

import (
	"log/slog"
	"memora/internal/config"
	"memora/internal/handlers/cards"
	"memora/internal/handlers/decks"
	"memora/internal/handlers/docs"
	"memora/internal/handlers/status"
	"memora/internal/handlers/users"
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
		userRoute := v1.Group("/users")
		{
			userRoute.GET("/:id", users.GetUser(services.Users))
			userRoute.POST("/", users.CreateUser(services.Users))
			userRoute.PATCH("/:id", users.PatchUser(services.Users))
			userRoute.DELETE("/:id", users.DeleteUser(services.Users))
		}
		cardRoute := v1.Group("/cards")
		{
			cardRoute.GET("/:id", cards.GetCard(services.Cards))
			cardRoute.POST("/", cards.CreateCard(services.Cards))
			cardRoute.PUT("/:id", cards.PutCard(services.Cards))
			cardRoute.DELETE("/:id", cards.DeleteCard(services.Cards))
		}
		deckRoute := v1.Group("/decks")
		{
			deckRoute.GET("/:id", decks.GetDeck(services.Decks))
			deckRoute.POST("/", decks.CreateDeck(services.Decks))
			deckRoute.PATCH("/:id", decks.PatchDeck(services.Decks))
			deckRoute.DELETE("/:id", decks.DeleteDeck(services.Decks))
		}
	}
}
