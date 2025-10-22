package users

import (
	"memora/internal/errors"
	"memora/internal/models"
	"memora/internal/services"
	"memora/internal/utils"
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
		id, err := utils.GetUID(c)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			return
		}

		user, err := userRepo.GetUser(c.Request.Context(), id)
		if errors.HandleError(c, err) {
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

// @Summary GET a users' owned decks from firestore by their ID
// @Description Return the user's owned decks
// @Tags Users
// @Produce json
// @Success 200 {object} []models.DisplayDeck
// @Router /api/v1/users/{id}/decks/owned [get]
// Return the users' owned decks based on an id
func GetDecks(userRepo *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := utils.GetUID(c)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			return
		}

		decks, err := userRepo.GetDecks(c.Request.Context(), id)
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
		uid := c.GetString("uid")

		if err := c.ShouldBindBodyWithJSON(&content); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid JSON",
			})
			return
		}

		content.Email = c.GetString("email")

		if err := userRepo.RegisterNewUser(c.Request.Context(), content, uid); err != nil {
			errors.HandleError(c, err)
			return
		}

		c.JSON(http.StatusCreated, models.ReturnID{
			ID: uid,
		})
	}
}

// @Summary Update a user in firestore by their ID
// @Description Return updated user information
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.PatchUser true "User info to update"
// @Success 200 {object} models.User
// @Router /api/v1/users/{id} [patch]
// Updates a user based on an id and returns the updated user
func PatchUser(userRepo *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var updates models.CreateUser
		id := c.GetString("uid")
		if err := c.ShouldBindBodyWithJSON(&updates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid JSON",
			})
			return
		}

		user, err := userRepo.UpdateUser(c.Request.Context(), updates, id)
		if errors.HandleError(c, err) {
			return
		}

		c.JSON(http.StatusOK, user)
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
		id := c.GetString("uid")

		err := userRepo.DeleteUser(c.Request.Context(), id)
		if errors.HandleError(c, err) {
			return
		}

		c.Status(http.StatusNoContent)
	}
}
