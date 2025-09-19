package decks

import (
	"errors"
	"memora/internal/models"
	"memora/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateDeck(deckRepo *services.DeckService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var content models.CreateDeck

		if err := c.ShouldBindBodyWithJSON(&content); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid body",
			})
			return
		}

		id, err := deckRepo.RegisterNewDeck(c.Request.Context(), content)
		if err != nil {
			switch {
			case errors.Is(err, models.ErrUserNotFound):
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "user does not exist",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})
			}
			return
		}
		c.JSON(http.StatusOK, models.ReturnID{
			ID: id,
		})
	}
}
