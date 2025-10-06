package errors

import (
	"errors"
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
	ErrorMap                  = map[error]struct {
		Status  int
		Message string
	}{
		ErrNotFound:               {Status: http.StatusBadRequest, Message: "did not find document"},
		ErrInvalidUser:            {Status: http.StatusBadRequest, Message: "invalid user data"},
		ErrInvalidCard:            {Status: http.StatusBadRequest, Message: "invalid card data"},
		ErrInvalidId:              {Status: http.StatusBadRequest, Message: "invalid id"},
		ErrInvalidDeck:            {Status: http.StatusBadRequest, Message: "invalid deck, missing fields"},
		ErrInvalidEmailNotPresent: {Status: http.StatusBadRequest, Message: "email not registerd"},
		ErrInvalidEmailPresent: {Status: http.StatusBadRequest, Message: "email alredy registerd"},
	}
)

func HandleError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	for k, v := range ErrorMap {
		if errors.Is(err, k) {
			c.JSON(v.Status, gin.H{"error": v.Message})
			return true
		}
	}

	// fallback for unexpected errors
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	return true
}
