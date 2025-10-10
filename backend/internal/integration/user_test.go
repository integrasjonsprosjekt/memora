package integration_test

import (
	"strings"
	"testing"
)

func TestCreateUser(t *testing.T) {

	r := SetupRouter(t)

	t.Run("Missing parameter in body", func(t *testing.T) {
		body := `{
			"email": "herman@example.com",
			"password": "short"
		}`
		w := PerformRequest(r, "POST", "/api/v1/users/", strings.NewReader(body))
		if w.Code != 400 {
			t.Errorf("Expected status code 400, got %d", w.Code)
		}
		resp := w.Body.String()
		expected := `{"error":"invalid user data"}`
		if resp != expected {
			t.Errorf("Expected response body %q, got %q", expected, resp)
		}
	})

	t.Run("Invalid email", func(t *testing.T) {
		body := `{
			"email": "not-an-email",
			"password": "validpassword",
			"name": "Herman"
		}`
		w := PerformRequest(r, "POST", "/api/v1/users/", strings.NewReader(body))
		if w.Code != 400 {
			t.Errorf("Expected status code 400, got %d", w.Code)
		}
		resp := w.Body.String()
		expected := `{"error":"invalid user data"}`
		if resp != expected {
			t.Errorf("Expected response body %q, got %q", expected, resp)
		}
	})

	t.Run("Valid user creation", func(t *testing.T) {
		body := `{
			"email": "herman@example.com",
			"password": "validpassword",
			"name": "Herman"
		}`
		w := PerformRequest(r, "POST", "/api/v1/users/", strings.NewReader(body))
		if w.Code != 201 {
			t.Errorf("Expected status code 201, got %d", w.Code)
		}
		resp := w.Body.String()
		expectedPrefix := `{"id":"`
		if !strings.HasPrefix(resp, expectedPrefix) {
			t.Errorf("Expected response body to start with %q, got %q", expectedPrefix, resp)
		}
	})

	t.Run("Duplicate email", func(t *testing.T) {
		body := `{
			"email": "herman@example.com",
			"password": "anotherpassword",
			"name": "Herman2"
		}`
		w := PerformRequest(r, "POST", "/api/v1/users/", strings.NewReader(body))
		if w.Code != 400 {
			t.Errorf("Expected status code 400, got %d", w.Code)
		}
		resp := w.Body.String()
		expected := `{"error":"email already registered"}`
		if resp != expected {
			t.Errorf("Expected response body %q, got %q", expected, resp)
		}
	})
}
