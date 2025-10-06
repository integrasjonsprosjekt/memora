package decks

import (
	"memora/internal/errors"
	"memora/internal/models"
	"memora/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateDeck(deckRepo *services.DeckService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var content models.CreateDeck
		content.SharedEmails = []string{}

		if err := c.ShouldBindBodyWithJSON(&content); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid body",
			})
			return
		}

		id, err := deckRepo.RegisterNewDeck(c.Request.Context(), content)
		if errors.HandleError(c, err) {
			return
		}
		c.JSON(http.StatusOK, models.ReturnID{
			ID: id,
		})
	}
}

func GetDeck(deckRepo *services.DeckService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		deck, err := deckRepo.GetOneDeck(c.Request.Context(), id)
		if errors.HandleError(c, err) {
			return
		}
		c.JSON(http.StatusOK, deck)
	}
}

func PatchDeck(deckRepo *services.DeckService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body models.UpdateDeck
		id := c.Param("id")

		if err := c.ShouldBindBodyWithJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid body",
			})
			return
		}

		deck, err := deckRepo.UpdateDeck(c.Request.Context(), id, body)
		if errors.HandleError(c, err) {
			return
		}

		c.JSON(http.StatusOK, deck)
	}
}
