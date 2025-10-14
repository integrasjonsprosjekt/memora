package middleware

import (
	"context"
	"memora/internal/errors"
	"strings"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
)

func FirebaseAuthMiddleware(app *firebase.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		authClient, err := app.Auth(context.Background())
		if err != nil {
			errors.HandleError(c, err)
			c.Abort()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			errors.HandleError(c, errors.ErrUnauthorized)
			c.Abort()
			return
		}

		idToken := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := authClient.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			errors.HandleError(c, errors.ErrUnauthorized)
			c.Abort()
			return
		}

		c.Set("uid", token.UID)
		c.Next()
	}
}
