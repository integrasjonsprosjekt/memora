package decks

import (
	"memora/internal/errors"
	"memora/internal/models"
	"memora/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Get a deck
// @Description Retrieves card information from Firestore by its ID
// @Tags Decks
// @Accept json
// @Produce json
// @Param id path string true "Deck ID"
// @Success 200 {object} models.DeckResponse
// @Router /api/v1/decks/{id} [get]
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

func GetCardInDeck(deckRepo *services.DeckService) gin.HandlerFunc {
	return func(c *gin.Context) {
		deckID := c.Param("deckID")
		cardID := c.Param("cardID")

		card, err := deckRepo.GetCardInDeck(c.Request.Context(), deckID, cardID)
		if errors.HandleError(c, err) {
			return
		}
		c.JSON(http.StatusOK, card)
	}
}

// @Summary Create a deck
// @Description Creates a new deck in Firestore and returns its ID
// @Tags Decks
// @Accept json
// @Produce json
// @Param user body models.CreateDeck true "Deck info"
// @Success 201 {object} models.ReturnID
// @Router /api/v1/decks [post]
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
		c.JSON(http.StatusCreated, models.ReturnID{
			ID: id,
		})
	}
}

func CreateCardInDeck(deckRepo *services.DeckService) gin.HandlerFunc {
	return func(c *gin.Context) {
		deckID := c.Param("deckID")

		rawData, err := c.GetRawData()
		if errors.HandleError(c, err) {
			return
		}

		deck, err := deckRepo.AddCardToDeck(c.Request.Context(), deckID, rawData)
		if errors.HandleError(c, err) {
			return
		}

		c.JSON(http.StatusCreated, deck)
	}
}

// @Summary Update a deck
// @Description Updates a decks data in Firestore by ID
// @Tags Decks
// @Accept json
// @Produce json
// @Param user body models.UpdateDeck true "Deck info"
// @Success 200 {object} models.DeckResponse
// @Router /api/v1/decks/{id} [patch]
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

// @Summary Update a decks' emails
// @Description Updates a decks shared emails in Firestore by ID
// @Tags Decks
// @Accept json
// @Produce json
// @Param user body models.UpdateDeckEmails true "Deck info"
// @Success 200 {object} models.DeckResponse
// @Router /api/v1/decks/{id}/emails [patch]
func UpdateEmails(deckRepo *services.DeckService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var body models.UpdateDeckEmails

		if err := c.ShouldBindBodyWithJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid body",
			})
			return
		}

		deck, err := deckRepo.UpdateEmailsInDeck(c.Request.Context(), id, body)
		if errors.HandleError(c, err) {
			return
		}

		c.JSON(http.StatusOK, deck)
	}
}

// @Summary Update a decks' cards
// @Description Updates a decks cards in Firestore by ID
// @Tags Decks
// @Accept json
// @Produce json
// @Param user body models.UpdateDeckCards true "Deck info"
// @Success 200 {object} models.DeckResponse
// @Router /api/v1/decks/{id}/cards [patch]
func UpdateCards(deckRepo *services.DeckService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var body models.UpdateDeckCards

		if err := c.ShouldBindBodyWithJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid body",
			})
			return
		}

		deck, err := deckRepo.UpdateCardsInDeck(c.Request.Context(), id, body)
		if errors.HandleError(c, err) {
			return
		}

		c.JSON(http.StatusOK, deck)
	}
}

// @Summary Delete a deck
// @Description Deletes a deck from Firestore by its ID
// @Tags Decks
// @Accept json
// @Produce json
// @Param id path string true "Deck ID"
// @Success 204
// @Router /api/v1/decks/{id} [delete]
func DeleteDeck(deckRepo *services.DeckService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		err := deckRepo.DeleteDeck(c.Request.Context(), id)
		if errors.HandleError(c, err) {
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func DeleteCardInDeck(deckRepo *services.DeckService) gin.HandlerFunc {
	return func(c *gin.Context) {
		deckID := c.Param("deckID")
		cardID := c.Param("cardID")

		err := deckRepo.DeleteCardInDeck(c.Request.Context(), deckID, cardID)
		if errors.HandleError(c, err) {
			return
		}
		c.Status(http.StatusNoContent)
	}
}
