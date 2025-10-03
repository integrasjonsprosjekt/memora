package cards

import (
	"memora/internal/errors"
	"memora/internal/models"
	"memora/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Create a card
// @Description Creates a new card in Firestore and returns its ID
// @Tags Cards
// @Accept json
// @Produce json
// @Param card body object true "Card info (can be MultipleChoiceCard, FrontBackCard, OrderedCard, or BlanksCard)"
// @Success 201 {object} models.ReturnID
// @Router /api/v1/cards [post]
func CreateCard(cardRepo *services.CardService) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawData, err := c.GetRawData()
		if errors.HandleError(c, err) {
			return
		}

		id, err := cardRepo.CreateCard(c.Request.Context(), rawData)
		if errors.HandleError(c, err) {
			return
		}

		c.JSON(http.StatusCreated, models.ReturnID{
			ID: id,
		})
	}
}

// @Summary Get a card
// @Description Retrieves card information from Firestore by its ID
// @Tags Cards
// @Accept json
// @Produce json
// @Param id path string true "Card ID"
// @Success 200 {object} models.AnyCard
// @Router /api/v1/cards/{id} [get]
func GetCard(cardRepo *services.CardService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		card, err := cardRepo.GetCard(c.Request.Context(), id)
		if errors.HandleError(c, err) {
			return
		}

		c.JSON(http.StatusOK, card)
	}
}

// @Summary Update a card
// @Description Updates card data in Firestore by ID
// @Tags Cards
// @Accept json
// @Produce json
// @Param id path string true "Card ID"
// @Success 200 {object} models.AnyCard
// @Router /api/v1/cards/{id} [patch]
func PutCard(cardRepo *services.CardService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		rawData, err := c.GetRawData()
		if errors.HandleError(c, err) {
			return
		}

		card, err := cardRepo.UpdateCard(c.Request.Context(), rawData, id)
		if errors.HandleError(c, err) {
			return
		}

		c.JSON(http.StatusOK, card)
	}
}

// @Summary Delete a card
// @Description Deletes a card from Firestore by its ID
// @Tags Cards
// @Accept json
// @Produce json
// @Param id path string true "Card ID"
// @Success 204
// @Router /api/v1/cards/{id} [delete]
func DeleteCard(cardRepo *services.CardService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		err := cardRepo.DeleteCard(c.Request.Context(), id)
		if errors.HandleError(c, err) {
			return
		}

		c.Status(http.StatusNoContent)
	}
}
