package users

import (
	"errors"
	"memora/internal/models"
	"memora/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary GET a user from firestore by their ID
// @Description Returns user information
// @Tags Users
// @Produce json
// @Success 200 {object} models.User
// @Router /api/v1/users/{id} [get]
// Returns the user based on an id
func GetUser(userRepo *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := userRepo.GetUser(c.Request.Context(), c.Param("id"))
		if err != nil {
			if errors.Is(err, models.ErrUserNotFound) {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "does not exist",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

// @Summary Create a user and returns their ID
// @Description Creates a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.CreateUser true "User info"
// @Success 201 {object} ReturnID
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
		if err != nil {
			if errors.Is(err, models.ErrInvalidUser) {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "invalid JSON",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			return
		}

		c.JSON(http.StatusCreated, ReturnID{
			ID: id,
		})
	}
}
