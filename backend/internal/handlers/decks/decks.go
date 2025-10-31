package decks

import (
	"context"
	"memora/internal/errors"
	"memora/internal/models"
	"memora/internal/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary Get cards in a deck
// @Description Retrieves cards from a specified deck in Firestore
// @Tags Decks
// @Accept json
// @Produce json
// @Param deckID path string true "Deck ID"
// @Param limit query string false "Number of cards to retrieve" default(20)
// @Param cursor query string false "Cursor for pagination"
// @Success 200 {object} models.CardsResponse
// @Router /api/v1/decks/{deckID}/cards [get]
func GetCardsInDeck(deckRepo *services.DeckService) gin.HandlerFunc {
	return func(c *gin.Context) {
		deckID := c.Param("deckID")
		limit := c.DefaultQuery("limit", "20")
		cursor := c.DefaultQuery("cursor", "")

		cards, hasMore, err := deckRepo.GetCardsInDeck(c.Request.Context(), deckID, limit, cursor)
		if errors.HandleError(c, err) {
			return
		}
		c.JSON(http.StatusOK, models.CardsResponse{
			Cards:   cards,
			HasMore: hasMore,
		})
	}
}

// @Summary Get a deck
// @Description Retrieves deck information from Firestore by its ID
// @Tags Decks
// @Accept json
// @Produce json
// @Param deckID path string true "Deck ID"
// @Success 200 {object} models.DeckResponse
// @Router /api/v1/decks/{deckID} [get]
func GetDeck(deckRepo *services.DeckService) gin.HandlerFunc {
	return func(c *gin.Context) {
		deckID := c.Param("deckID")
		filter := c.DefaultQuery("filter", "title,owner_id,shared_emails")

		deck, err := deckRepo.GetOneDeck(c.Request.Context(), deckID, filter)
		if errors.HandleError(c, err) {
			return
		}
		c.JSON(http.StatusOK, deck)
	}
}

// @Summary Get a card in a deck
// @Description Retrieves card information from Firestore by its ID
// @Tags Decks
// @Accept json
// @Produce json
// @Param deckID path string true "Deck ID"
// @Param cardID path string true "Card ID"
// @Success 200 {object} models.AnyCard
// @Router /api/v1/decks/{deckID}/cards/{cardID} [get]
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

// @Summary Create a card in a deck
// @Description Creates a new card in a specified deck and returns the cards
// @Tags Decks
// @Accept json
// @Produce json
// @Param deckID path string true "Deck ID"
// @Param card body object true "Card info (can be MultipleChoiceCard, FrontBackCard, OrderedCard, or BlanksCard)"
// @Success 201 {object} []models.Card
// @Router /api/v1/decks/{deckID}/cards [post]
func CreateCardInDeck(deckRepo *services.DeckService) gin.HandlerFunc {
	return func(c *gin.Context) {
		deckID := c.Param("deckID")

		rawData, err := c.GetRawData()
		if errors.HandleError(c, err) {
			return
		}

		cards, err := deckRepo.AddCardToDeck(c.Request.Context(), deckID, rawData)
		if errors.HandleError(c, err) {
			return
		}

		c.JSON(http.StatusCreated, cards)
	}
}

// @Summary Update a deck
// @Description Updates a decks data in Firestore by ID
// @Tags Decks
// @Accept json
// @Produce json
// @Param user body models.UpdateDeck true "Deck info"
// @Success 200 {object} models.DeckResponse
// @Router /api/v1/decks/{deckID} [patch]
func PatchDeck(deckRepo *services.DeckService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body models.UpdateDeck
		id := c.Param("deckID")

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
// @Router /api/v1/decks/{deckID}/emails [patch]
func UpdateEmails(deckRepo *services.DeckService) gin.HandlerFunc {
	return func(c *gin.Context) {
		deckID := c.Param("deckID")

		var body models.UpdateDeckEmails

		if err := c.ShouldBindBodyWithJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid body",
			})
			return
		}

		deck, err := deckRepo.UpdateEmailsInDeck(c.Request.Context(), deckID, body)
		if errors.HandleError(c, err) {
			return
		}

		c.JSON(http.StatusOK, deck)
	}
}

// @Summary Update a card in a deck
// @Description Updates a card in a deck in Firestore by ID
// @Tags Decks
// @Accept json
// @Produce json
// @Param user body models.AnyCard true "Deck info"
// @Param deckID path string true "Deck ID"
// @Param cardID path string true "Card ID"
// @Success 200 {object} []models.Card
// @Router /api/v1/decks/{deckID}/cards/{cardID} [put]
func UpdateCard(deckRepo *services.DeckService) gin.HandlerFunc {
	return func(c *gin.Context) {
		deckID := c.Param("deckID")
		cardID := c.Param("cardID")

		rawData, err := c.GetRawData()
		if errors.HandleError(c, err) {
			return
		}

		card, err := deckRepo.UpdateCardInDeck(
			c.Request.Context(),
			deckID, cardID,
			rawData,
		)

		if errors.HandleError(c, err) {
			return
		}

		c.JSON(http.StatusOK, card)
	}
}

// @Summary Delete a deck along with its cards
// @Description Deletes a deck from Firestore by its ID
// @Tags Decks
// @Accept json
// @Produce json
// @Param deckID path string true "Deck ID"
// @Success 204
// @Router /api/v1/decks/{deckID} [delete]
func DeleteDeck(deckRepo *services.DeckService) gin.HandlerFunc {
	return func(c *gin.Context) {
		deckID := c.Param("deckID")

		err := deckRepo.DeleteDeck(c.Request.Context(), deckID)
		if errors.HandleError(c, err) {
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// @Summary Delete a card in a deck
// @Description Deletes a card from a specified deck by its ID
// @Tags Decks
// @Accept json
// @Produce json
// @Param deckID path string true "Deck ID"
// @Param cardID path string true "Card ID"
// @Success 204
// @Router /api/v1/decks/{deckID}/cards/{cardID} [delete]
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

func GetDueCardsInDeck(deckRepo *services.DeckService) gin.HandlerFunc {
	return func(c *gin.Context) {
		deckID := c.Param("deckID")
		userID := c.Param("userID")
		limit := c.DefaultQuery("limit", "20")
		cursor := c.DefaultQuery("cursor", "")

		limitInt, err := strconv.Atoi(limit)
		if errors.HandleError(c, err) {
			return
		}

		cards, nextCursor, hasMore, err := deckRepo.GetDueCardsInDeck(c.Request.Context(), deckID, userID, limitInt, cursor)
		if errors.HandleError(c, err) {
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"cards":       cards,
			"next_cursor": nextCursor,
			"has_more":    hasMore,
		})
	}
}

func GetProgress(deckRepo *services.DeckService) gin.HandlerFunc {
	return func(c *gin.Context) {
		deckID := c.Param("deckID")
		cardID := c.Param("cardID")
		userID := c.Param("userID")

		progress, err := deckRepo.GetCardProgress(c.Request.Context(), deckID, cardID, userID)
		if errors.HandleError(c, err) {
			return
		}

		c.JSON(http.StatusOK, progress)
	}
}

func UpdateProgress(deckRepo *services.DeckService) gin.HandlerFunc {
	return func(c *gin.Context) {
		deckID := c.Param("deckID")
		cardID := c.Param("cardID")
		userID := c.Param("userID")

		var body models.CardRating

		if err := c.ShouldBindBodyWithJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid body",
			})
			return
		}

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

			defer cancel()

			if err := deckRepo.UpdateCardProgress(ctx, deckID, cardID, userID, body); err != nil {
				errors.HandleError(c, err)
			}
		}()

		c.Status(http.StatusAccepted)
	}
}
