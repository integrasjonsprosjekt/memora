package errors

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrNotFound               = errors.New("could not find id")
	ErrInvalidUser            = errors.New("invalid user data")
	ErrInvalidCard            = errors.New("invalid card data")
	ErrInvalidDeck            = errors.New("invalid deck data")
	ErrInvalidEmailNotPresent = errors.New("email not registerd")
	ErrInvalidEmailPresent    = errors.New("email alredy registerd")
	ErrInvalidId              = errors.New("invalid id")
	ErrFailedUpdatingEmail    = errors.New("failed to update emails")
	ErrFailedUpdatingCards    = errors.New("failed to update cards")
	ErrAlreadyExists          = errors.New("resource already exists")
	ErrUnauthorized           = errors.New("unauthorized")
	ErrorMap                  = map[error]struct {
		Status  int
		Message string
	}{
		ErrNotFound: {
			Status:  http.StatusBadRequest,
			Message: "did not find document",
		},
		ErrInvalidUser: {Status: http.StatusBadRequest, Message: "invalid user data"},
		ErrInvalidCard: {Status: http.StatusBadRequest, Message: "invalid card data"},
		ErrInvalidId:   {Status: http.StatusBadRequest, Message: "invalid id"},
		ErrInvalidDeck: {
			Status:  http.StatusBadRequest,
			Message: "invalid deck, missing fields",
		},
		ErrInvalidEmailNotPresent: {Status: http.StatusBadRequest, Message: "email not registered"},
		ErrInvalidEmailPresent: {
			Status:  http.StatusBadRequest,
			Message: "email already registered",
		},
		ErrFailedUpdatingEmail: {
			Status:  http.StatusBadRequest,
			Message: "failed to update emails",
		},
		ErrFailedUpdatingCards: {
			Status:  http.StatusBadRequest,
			Message: "failed to update cards",
		},
		ErrAlreadyExists: {
			Status:  http.StatusConflict,
			Message: "resource already exists",
		},
		ErrUnauthorized: {Status: http.StatusUnauthorized, Message: "unauthorized operation"},
	}
)

func HandleError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	for k, v := range ErrorMap {
		if errors.Is(err, k) {
			log.Println("Got error: ", err)
			c.JSON(v.Status, gin.H{"error": v.Message})
			return true
		}
	}

	// fallback for unexpected errors
	log.Println("Got unexpected error: ", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	return true
}
