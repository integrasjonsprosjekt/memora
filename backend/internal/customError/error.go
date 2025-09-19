package customerror

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidUser  = errors.New("invalid user data")
	ErrorMap        = map[error]struct {
		Status  int
		Message string
	}{
		ErrUserNotFound: {Status: http.StatusBadRequest, Message: "user not found"},
		ErrInvalidUser:  {Status: http.StatusBadRequest, Message: "invalid user data"},
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
