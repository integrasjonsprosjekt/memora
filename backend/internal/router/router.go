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
			userRoute.GET(
				"/",
				middleware.FirebaseAuthMiddleware(services.Auth),
				users.GetUser(services.Users),
			)
			userRoute.POST(
				"/",
				middleware.FirebaseAuthMiddleware(services.Auth),
				users.CreateUser(services.Users),
			)
			userRoute.PATCH(
				"/",
				middleware.FirebaseAuthMiddleware(services.Auth),
				users.PatchUser(services.Users),
			)
			userRoute.DELETE(
				"/",
				middleware.FirebaseAuthMiddleware(services.Auth),
				users.DeleteUser(services.Users),
			)
			userRoute.GET(
				"/decks",
				middleware.FirebaseAuthMiddleware(services.Auth),
				users.GetDecks(services.Users),
			)
		}

		// Deck-related endpoints
		deckRoute := v1.Group("/decks")
		{
			deckRoute.GET(
				"/:deckID",
				middleware.FirebaseAuthMiddleware(services.Auth),
				decks.GetDeck(services.Decks),
			)
			deckRoute.POST(
				"/",
				middleware.FirebaseAuthMiddleware(services.Auth),
				decks.CreateDeck(services.Decks),
			)
			deckRoute.DELETE(
				"/:deckID",
				middleware.FirebaseAuthMiddleware(services.Auth),
				decks.DeleteDeck(services.Decks),
			)

			// PATCH routes for updating specific fields of a deck
			deckRoute.PATCH(
				"/:deckID",
				middleware.FirebaseAuthMiddleware(services.Auth),
				decks.PatchDeck(services.Decks),
			)
			deckRoute.PATCH(
				"/:deckID/emails",
				middleware.FirebaseAuthMiddleware(services.Auth),
				decks.UpdateEmails(services.Decks),
			)

			cardRoute := deckRoute.Group("/:deckID/cards")
			{
				cardRoute.GET(
					"/due",
					middleware.FirebaseAuthMiddleware(services.Auth),
					decks.GetDueCardsInDeck(services.Decks),
				)
				cardRoute.GET(
					"/",
					middleware.FirebaseAuthMiddleware(services.Auth),
					decks.GetCardsInDeck(services.Decks),
				)
				cardRoute.POST(
					"/",
					middleware.FirebaseAuthMiddleware(services.Auth),
					decks.CreateCardInDeck(services.Decks),
				)
				cardRoute.GET(
					"/:cardID",
					middleware.FirebaseAuthMiddleware(services.Auth),
					decks.GetCardInDeck(services.Decks),
				)
				cardRoute.PUT(
					"/:cardID",
					middleware.FirebaseAuthMiddleware(services.Auth),
					decks.UpdateCard(services.Decks),
				)
				cardRoute.DELETE(
					"/:cardID",
					middleware.FirebaseAuthMiddleware(services.Auth),
					decks.DeleteCardInDeck(services.Decks),
				)
				progress := cardRoute.Group("/:cardID/progress")
				{
					progress.GET(
						"/",
						middleware.FirebaseAuthMiddleware(services.Auth),
						decks.GetProgress(services.Decks),
					)
					progress.PUT(
						"/",
						middleware.FirebaseAuthMiddleware(services.Auth),
						decks.UpdateProgress(services.Decks),
					)
				}
			}
		}
	}
}
