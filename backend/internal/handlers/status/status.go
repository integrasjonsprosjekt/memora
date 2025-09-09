package status

import (
	"memora/internal/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary Status page for the API
// @Description Returns version and uptime
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Router api/v1/status [get]
// Returns status of the service, and the services used by it
func GetStatus(c *gin.Context) {
	status := Status{
		Version: "v1",
		Uptime:  time.Duration(time.Since(config.StartTime).Seconds()),
	}

	c.JSON(http.StatusOK, status)
}
