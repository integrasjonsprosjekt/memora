package docs

import (
	"net/http"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/gin-gonic/gin"
)

func GetDocs(c *gin.Context) {
	htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
		SpecURL: "./docs/swagger.json",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "Memora",
		},
		DarkMode: true,
	})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Send HTML to client
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlContent))
}
