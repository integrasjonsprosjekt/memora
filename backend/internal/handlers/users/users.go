package users

import (
	"memora/internal/errors"
	"memora/internal/models"
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

// @Summary GET a users' owned decks from firestore by their ID
// @Description Return the user's owned decks
// // @Tags Users
// @Produce json
// @Success 200 {object} []models.DisplayDeck
// @Router /api/v1/users/{id}/decks/owned [get]
// Return the users' owned decks based on an id
func GetDecksOwned(userRepo *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		decks, err := userRepo.GetDecksOwned(c.Request.Context(), c.Param("id"))
		if errors.HandleError(c, err) {
			return
		}

		c.JSON(http.StatusOK, decks)
	}
}

// @Summary GET a users' shared decks from firestore by their ID
// @Description Return the user's shared decks
// // @Tags Users
// @Produce json
// @Success 200 {object} []models.DisplayDeck
// @Router /api/v1/users/{id}/decks/shared [get]
// Return the users' shared decks based on an id
func GetDecksShared(userRepo *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		decks, err := userRepo.GetDecksShared(c.Request.Context(), c.Param("id"))
		if errors.HandleError(c, err) {
			return
		}

		c.JSON(http.StatusOK, decks)
	}
}

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

// @Summary Deletes a user from firestore by their ID
// @Description Return card information
// @Tags Users
// @Accept json
// @Produce json
// @Success 204
// @Router /api/v1/users/{id} [delete]
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
