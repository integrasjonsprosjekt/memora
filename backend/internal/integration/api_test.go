package integration_test

import (
	"log"
	"memora/internal/firebase"
	"memora/internal/router"
	"memora/internal/services"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func SetupRouter(t *testing.T) *gin.Engine {
	client, err := firebase.InitEmulator()
	if err != nil {
		log.Fatal(err)
	}

	validate := validator.New()
	repos := firebase.NewRepositories(client)
	svc := services.NewServices(repos, validate)

	r := router.New()
	router.Route(r, svc)

	return r
}

func PerformRequest(r *gin.Engine, method, path string, body *strings.Reader) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	req.Header.Set("content-type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
