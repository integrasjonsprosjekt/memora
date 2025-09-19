package decks

import (
	customerror "memora/internal/customError"
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
		if customerror.HandleError(c, err) {
			return
		}
		c.JSON(http.StatusOK, models.ReturnID{
			ID: id,
		})
	}
}

func GetDeck(deckRepo *services.DeckService) gin.HandlerFunc {
	return func(c *gin.Context) {
		deck, err := deckRepo.GetOneDeck(c.Request.Context(), c.Param("id"))
		if customerror.HandleError(c, err) {
			return
		}
		c.JSON(http.StatusOK, deck)
	}
}
