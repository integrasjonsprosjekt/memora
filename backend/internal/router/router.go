package router

import (
	"log/slog"
	"memora/internal/config"
	"memora/internal/handlers/decks"
	"memora/internal/handlers/docs"
	"memora/internal/handlers/status"
	"memora/internal/handlers/users"
	"memora/internal/middleware"
	"memora/internal/services"

	"github.com/gin-gonic/gin"
)

// New initializes a new Gin router with middleware.
func New() *gin.Engine {
	router := gin.New()

	logger := slog.Default()

	router.Use(gin.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.Logging(logger))

	return router
}

// Route sets up all the routes for the application.
func Route(c *gin.Engine, services *services.Services) {
	// API v1 routes
	v1 := c.Group(config.BasePath)
	{
		// Health and status endpoints
		v1.GET("/status", status.GetStatus)
		v1.GET("/docs", docs.GetDocs)

		// User-related endpoints
		userRoute := v1.Group("/users")
		{
			userRoute.GET("/:id", users.GetUser(services.Users))
			userRoute.POST("/", users.CreateUser(services.Users))
			userRoute.PATCH("/:id", users.PatchUser(services.Users))
			userRoute.DELETE("/:id", users.DeleteUser(services.Users))

			userRoute.GET("/:id/decks", users.GetDecks(services.Users))
		}

		// Deck-related endpoints
		deckRoute := v1.Group("/decks")
		{
			deckRoute.GET("/:deckID", decks.GetDeck(services.Decks))
			deckRoute.POST("/", decks.CreateDeck(services.Decks))
			deckRoute.DELETE("/:deckID", decks.DeleteDeck(services.Decks))

			// PATCH routes for updating specific fields of a deck
			deckRoute.PATCH("/:deckID", decks.PatchDeck(services.Decks))
			deckRoute.PATCH("/:deckID/emails", decks.UpdateEmails(services.Decks))

			cardRoute := deckRoute.Group("/:deckID/cards")
			{
				cardRoute.GET("/", decks.GetCardsInDeck(services.Decks))
				cardRoute.POST("/", decks.CreateCardInDeck(services.Decks))
				cardRoute.GET("/:cardID", decks.GetCardInDeck(services.Decks))
				cardRoute.PUT("/:cardID", decks.UpdateCard(services.Decks))
				cardRoute.DELETE("/:cardID", decks.DeleteCardInDeck(services.Decks))
			}
			progress := deckRoute.Group("/:deckID/users/:userID/progress/:cardID")
			{
				progress.POST("/", decks.CreateProgress(services.Decks))
				progress.GET("/due", decks.GetDueCardsInDeck(services.Decks))
				progress.GET("/", decks.GetProgress(services.Decks))
				progress.PUT("/", decks.UpdateProgress(services.Decks))
			}
		}
	}
}
