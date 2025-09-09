package status

import (
	"memora/internal/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Returns status of the service, and the services used by it
func GetStatus(c *gin.Context) {
	status := Status{
		Version: "v1",
		Uptime:  time.Duration(time.Since(config.StartTime).Seconds()),
	}

	c.JSON(http.StatusOK, status)
}
