package users

import (
	"memora/internal/errors"
	"memora/internal/models"
	"memora/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Create a user and return their ID
// @Description Creates a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.CreateUser true "User info"
// @Success 200 {object} models.User
// @Router /api/v1/users [post]
// Creates a new user in firestore
func CreateUser(userRepo *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var content models.CreateUser

		if err := c.ShouldBindBodyWithJSON(&content); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid JSON",
			})
			return
		}

		id, err := userRepo.RegisterNewUser(c.Request.Context(), content)
		if errors.HandleError(c, err) {
			return
		}

		c.JSON(http.StatusCreated, models.ReturnID{
			ID: id,
		})
	}
}

// @Summary Patch the users' data by ID
// @Description Updates/replaces data
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.PatchUser true "User info"
// @Success 204
// @Router /api/v1/users/{id} [patch]
// Updates a user in firestore
func PatchUser(userRepo *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var updates models.PatchUser

		if err := c.ShouldBindBodyWithJSON(&updates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid JSON",
			})
			return
		}

		err := userRepo.UpdateUser(c.Request.Context(), updates, id)
		if errors.HandleError(c, err) {
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User updated succesfully",
		})
	}
}

func DeleteUser(userRepo *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		err := userRepo.DeleteUser(c.Request.Context(), id)
		if errors.HandleError(c, err) {
			return
		}

		c.Status(http.StatusNoContent)
	}
}
