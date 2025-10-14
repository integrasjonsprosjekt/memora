package middleware

import (
	"log"
	"memora/internal/errors"
	"memora/internal/services"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func FirebaseAuthMiddleware(auth *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			errors.HandleError(c, errors.ErrUnauthorized)
			c.Abort()
			return
		}

		idToken := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := auth.VerifyIDToken(c.Request.Context(), idToken)
		if err != nil {
			errors.HandleError(c, err)
			c.Abort()
			return
		}

		c.Set("uid", token.UID)
		c.Set("email", token.Claims["email"])
		c.Next()
	}
}
