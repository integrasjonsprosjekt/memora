package cards

import (
	"memora/internal/errors"
	"memora/internal/models"
	"memora/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
