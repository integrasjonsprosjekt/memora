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
	client, app, err := firebase.InitEmulator()
	if err != nil {
		log.Fatal(err)
	}

	auth, err := firebase.NewFirebaseAuth(app)
	if err != nil {
		log.Fatal(err)
	}

	validate := validator.New()
	repos := firebase.NewRepositories(client, auth)
	svc := services.NewServices(repos, validate)

	r := router.New()
	router.Route(r, svc)

	return r
}

func PerformRequest(
	r *gin.Engine,
	method, path string,
	body io.Reader,
	token string,
) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	req.Header.Set("content-type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func CreateTestUser1(t *testing.T) string {
	url := "http://127.0.0.1:9099/identitytoolkit.googleapis.com/v1/accounts:signUp?key=any"
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

	defer func() { _ = resp.Body.Close() }()

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}

	idToken, ok := result["idToken"].(string)
	if !ok || idToken == "" {
		t.Fatal("failed to get idToken from emulator")
	}

	return idToken
}

func CreateTestUser2(t *testing.T) string {
	url := "http://127.0.0.1:9099/identitytoolkit.googleapis.com/v1/accounts:signUp?key=any"
	payload := map[string]string{
		"email":             "test@user2.com",
		"password":          "verysecurepassword",
		"returnSecureToken": "true",
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}

	idToken, ok := result["idToken"].(string)
	if !ok || idToken == "" {
		t.Fatal("failed to get idToken from emulator")
	}

	return idToken
}
