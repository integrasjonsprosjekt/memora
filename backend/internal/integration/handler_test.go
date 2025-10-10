package integration_test

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestUser(t *testing.T) {

	// Store created user ID for further tests
	var userID, email string

	// Setup router
	r := SetupRouter(t)

	// POST /api/v1/users/ - Create a new user
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

		var respData struct {
			ID string `json:"id"`
		}

		if err := json.Unmarshal(w.Body.Bytes(), &respData); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if respData.ID == "" {
			t.Errorf("Expected a non-empty user ID, got empty string")
		}

		userID = respData.ID

		// Optional: check that the response contains only the "id" field
		var generic map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &generic); err != nil {
			t.Fatalf("Failed to unmarshal response for extra field check: %v", err)
		}
		if len(generic) != 1 {
			t.Errorf("Expected only 1 field in response, got %d", len(generic))
		}
	})

	// POST /api/v1/users/ - Create another user
	t.Run("Create another user", func(t *testing.T) {
		body := `{
			"email": "herman2@example.com",
			"password": "validpassword",
			"name": "Herman2"
		}`
		w := PerformRequest(r, "POST", "/api/v1/users/", strings.NewReader(body))
		if w.Code != 201 {
			t.Errorf("Expected status code 201, got %d", w.Code)
		}

		var respData struct {
			ID string `json:"id"`
		}

		if err := json.Unmarshal(w.Body.Bytes(), &respData); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if respData.ID == "" {
			t.Errorf("Expected a non-empty user ID, got empty string")
		}
		email = "herman2@example.com"
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

	// GET /api/v1/users/{userID} - Retrieve the created user
	t.Run("Get created user", func(t *testing.T) {
		if userID == "" {
			t.Fatal("userID is empty, cannot perform GET test")
		}
		path := "/api/v1/users/" + userID
		w := PerformRequest(r, "GET", path, nil)
		if w.Code != 200 {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}
		resp := w.Body.String()
		expectedSubstring := `"email":"herman@example.com"`
		if !strings.Contains(resp, expectedSubstring) {
			t.Errorf("Expected response body to contain %q, got %q", expectedSubstring, resp)
		}
	})

	// GET /api/v1/users/{invalidID} - Attempt to retrieve a non-existent user
	t.Run("Get non-existent user", func(t *testing.T) {
		invalidID := "nonexistentid"
		path := "/api/v1/users/" + invalidID
		w := PerformRequest(r, "GET", path, nil)
		if w.Code != 400 {
			t.Errorf("Expected status code 400, got %d", w.Code)
		}
		resp := w.Body.String()
		expected := `{"error":"did not find document"}`
		if resp != expected {
			t.Errorf("Expected response body %q, got %q", expected, resp)
		}
	})

	// PATCH /api/v1/users/{userID} - Update the created user
	t.Run("Update user name", func(t *testing.T) {
		if userID == "" {
			t.Fatal("userID is empty, cannot perform PATCH test")
		}
		body := `{
			"name": "Herman Updated"
		}`
		path := "/api/v1/users/" + userID
		w := PerformRequest(r, "PATCH", path, strings.NewReader(body))
		if w.Code != 200 {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}
		resp := w.Body.String()
		expectedSubstring := `"name":"Herman Updated"`
		if !strings.Contains(resp, expectedSubstring) {
			t.Errorf("Expected response body to contain %q, got %q", expectedSubstring, resp)
		}
	})

	// PATCH /api/v1/users/{invalidID} - Attempt to update a non-existent user
	t.Run("Update non-existent user", func(t *testing.T) {
		invalidID := "nonexistentid"
		body := `{
			"name": "Should Not Work"
		}`
		path := "/api/v1/users/" + invalidID
		w := PerformRequest(r, "PATCH", path, strings.NewReader(body))
		if w.Code != 400 {
			t.Errorf("Expected status code 400, got %d", w.Code)
		}
		resp := w.Body.String()
		expected := `{"error":"invalid id"}`
		if resp != expected {
			t.Errorf("Expected response body %q, got %q", expected, resp)
		}
	})

	// GET /api/v1/users/{userID} - Verify the update
	t.Run("Verify user update", func(t *testing.T) {
		if userID == "" {
			t.Fatal("userID is empty, cannot perform verification GET test")
		}
		path := "/api/v1/users/" + userID
		w := PerformRequest(r, "GET", path, nil)
		if w.Code != 200 {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}
		resp := w.Body.String()
		expectedSubstring := `"name":"Herman Updated"`
		if !strings.Contains(resp, expectedSubstring) {
			t.Errorf("Expected response body to contain %q, got %q", expectedSubstring, resp)
		}
	})

	// Decks and Shared tests can be added similarly once user tests are stable
	t.Run("Create one deck for user", func(t *testing.T) {
		if userID == "" {
			t.Fatal("userID is empty, cannot perform deck creation test")
		}
		body := `{
			"title": "Test Deck",
			"owner_id": "` + userID + `"
		}`
		w := PerformRequest(r, "POST", "/api/v1/decks/", strings.NewReader(body))
		if w.Code != 201 {
			t.Errorf("Expected status code 201, got %d", w.Code)
		}
	})

	t.Run("Update the created deck to share with another user", func(t *testing.T) {
		if userID == "" {
			t.Fatal("userID is empty, cannot perform deck update test")
		}
		body := `{
			"opp": "add",
			"shared_emails": ["` + email + `"]
		}`
		// First, we need to get the deck ID of the created deck
		w := PerformRequest(r, "GET", "/api/v1/users/"+userID+"/decks/owned", nil)
		if w.Code != 200 {
			t.Errorf("Expected status code 200 when fetching owned decks, got %d", w.Code)
			return
		}

		var decksResp []struct {
			ID string `json:"id"`
		}

		if err := json.Unmarshal(w.Body.Bytes(), &decksResp); err != nil {
			t.Fatalf("Failed to unmarshal owned decks response: %v", err)
			return
		}
		if len(decksResp) == 0 {
			t.Fatal("No owned decks found")
		}
		deckID := decksResp[0].ID

		// Now we can update the deck
		w = PerformRequest(r, "PATCH", "/api/v1/decks/"+deckID+"/emails", strings.NewReader(body))
		if w.Code != 200 {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}
		resp := w.Body.String()
		expectedSubstring := `"shared_emails":["` + email + `"]`
		if !strings.Contains(resp, expectedSubstring) {
			t.Errorf("Expected response body to contain %q, got %q", expectedSubstring, resp)
		}
	})
}
