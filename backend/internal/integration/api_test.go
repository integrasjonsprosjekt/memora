package integration_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"memora/internal/firebase"
	"memora/internal/router"
	"memora/internal/services"
	"net/http"
	"net/http/httptest"
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
	repos := firebase.NewRepositories(client, nil)
	svc := services.NewServices(repos, validate)

	r := router.New(svc.Auth)
	router.Route(r, svc)

	return r
}

func PerformRequest(r *gin.Engine, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	req.Header.Set("content-type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func CreateTestUser(r *gin.Engine, t *testing.T) string {
	url := "http://localhost:9099/identitytoolkit.googleapis.com/v1/accounts:signUp?key=any"
	payload := map[string]string{
		"email":             "test@user.com",
		"password":          "verysecurepassword",
		"returnSecureToken": "true",
	}
	body, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	t.Log(result)
	return result["idToken"].(string)
}
