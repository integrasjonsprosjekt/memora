package middleware

import (
	"memora/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimit(requestsPerMinute int, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUID(c)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Create a unique key for the user
		rateLimitKey := utils.UserKeyRateLimit(userID)
		ctx := c.Request.Context()

		// Use a pipeline to batch commands
		pipe := rdb.Pipeline()
		incr := pipe.Incr(ctx, rateLimitKey)
		pipe.Expire(ctx, rateLimitKey, time.Minute) // Set TTL to 1 minute
		_, err = pipe.Exec(ctx)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if incr.Val() > int64(requestsPerMinute) {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
	}
}
