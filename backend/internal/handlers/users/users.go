package users

import (
	"memora/internal/errors"
	"memora/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary GET a user from firestore by their ID
// @Description Return user information
// @Tags Users
// @Produce json
// @Success 200 {object} models.User
// @Router /api/v1/users/{id} [get]
// Return the user based on an id
func GetUser(userRepo *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := userRepo.GetUser(c.Request.Context(), c.Param("id"))
		if errors.HandleError(c, err) {
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
