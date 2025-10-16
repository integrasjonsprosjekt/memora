package integration_test

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestHandlers(t *testing.T) {

	// Store created user ID for further tests
	var deckID string
	email := "test@user2.com"

	// Setup router
	r := SetupRouter(t)

	token1 := CreateTestUser1(t)
	if token1 == "" {
		t.Fatal("Failed to create test user")
	}
	token2 := CreateTestUser2(t)
	if token2 == "" {
		t.Fatal("Failed to create second test user")
	}

	// POST /api/v1/users/ - Create a new user
	t.Run("Missing parameter in body", func(t *testing.T) {
		body := `{
			"test": "invalid",
		}`
		w := PerformRequest(r, "POST", "/api/v1/users/", strings.NewReader(body), token1)
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
        "name": "Herman"
    }`
		w := PerformRequest(r, "POST", "/api/v1/users/", strings.NewReader(body), token1)
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
			"name": "Herman2"
		}`
		w := PerformRequest(r, "POST", "/api/v1/users/", strings.NewReader(body), token2)
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
	})

	t.Run("Duplicate email", func(t *testing.T) {
		body := `{
			"name": "Herman2"
		}`
		w := PerformRequest(r, "POST", "/api/v1/users/", strings.NewReader(body), token2)
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
		w := PerformRequest(r, "GET", "/api/v1/users/", nil, token1)
		if w.Code != 200 {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}
		resp := w.Body.String()
		expectedSubstring := `"email":"test@user.com"`
		if !strings.Contains(resp, expectedSubstring) {
			t.Errorf("Expected response body to contain %q, got %q", expectedSubstring, resp)
		}
	})

	// GET /api/v1/users/{invalidID} - Attempt to retrieve a non-existent user
	t.Run("Get non-existent user without validation", func(t *testing.T) {
		invalidID := "nonexistentid"
		path := "/api/v1/users/" + invalidID
		w := PerformRequest(r, "GET", path, nil, "")
		if w.Code != 401 {
			t.Errorf("Expected status code 401, got %d", w.Code)
		}
		resp := w.Body.String()
		expected := `{"error":"unauthorized operation"}`
		if resp != expected {
			t.Errorf("Expected response body %q, got %q", expected, resp)
		}
	})

	// PATCH /api/v1/users/{userID} - Update the created user
	t.Run("Update user name", func(t *testing.T) {
		body := `{
			"name": "Herman Updated"
		}`
		w := PerformRequest(r, "PATCH", "/api/v1/users/", strings.NewReader(body), token1)
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
		w := PerformRequest(r, "PATCH", path, strings.NewReader(body), "")
		if w.Code != 401 {
			t.Errorf("Expected status code 401, got %d", w.Code)
		}
		resp := w.Body.String()
		expected := `{"error":"unauthorized operation"}`
		if resp != expected {
			t.Errorf("Expected response body %q, got %q", expected, resp)
		}
	})

	// GET /api/v1/users/{userID} - Verify the update
	t.Run("Verify user update", func(t *testing.T) {
		w := PerformRequest(r, "GET", "/api/v1/users/", nil, token1)
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
		body := `{
			"title": "Test Deck",
		}`
		w := PerformRequest(r, "POST", "/api/v1/decks/", strings.NewReader(body), token1)
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
			t.Errorf("Expected a non-empty deck ID, got empty string")
		}
		deckID = respData.ID
	})

	t.Run("Update the created deck to share with another user", func(t *testing.T) {
		body := `{
			"opp": "add",
			"shared_emails": ["` + email + `"]
		}`
		// Now we can update the deck
		w := PerformRequest(r, "PATCH", "/api/v1/decks/"+deckID+"/emails", strings.NewReader(body), token1)
		if w.Code != 200 {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}
		resp := w.Body.String()
		expectedSubstring := `"shared_emails":["` + email + `"]`
		if !strings.Contains(resp, expectedSubstring) {
			t.Errorf("Expected response body to contain %q, got %q", expectedSubstring, resp)
		}
	})

	t.Run("Add invalid email to deck", func(t *testing.T) {
		body := `{
			"opp": "add",
			"shared_emails": ["invalid-email"]
		}`
		// Now we can update the deck
		w := PerformRequest(r, "PATCH", "/api/v1/decks/"+deckID+"/emails", strings.NewReader(body), token1)
		if w.Code != 400 {
			t.Errorf("Expected status code 400, got %d", w.Code)
		}
		resp := w.Body.String()
		expected := `{"error":"email not registered"}`
		if resp != expected {
			t.Errorf("Expected response body %q, got %q", expected, resp)
		}
	})

	t.Run("Add cards to the created deck", func(t *testing.T) {
		body := `{
			"type": "front_back",
			"front": "What is the capital of France?",
			"back": "Paris"
		}`

		w := PerformRequest(r, "POST", "/api/v1/decks/"+deckID+"/cards/", strings.NewReader(body), token1)
		if w.Code != 201 {
			t.Errorf("Expected status code 201, got %d", w.Code)
		}

		resp := w.Body.String()
		expectedSubstring := `"front":"What is the capital of France?"`
		if !strings.Contains(resp, expectedSubstring) {
			t.Errorf("Expected response body to contain %q, got %q", expectedSubstring, resp)
		}
	})

	t.Run("Add multiple choice card to the created deck", func(t *testing.T) {
		body := `{
			"type": "multiple_choice",
			"question": "What is 2 + 2?",
			"options": {
				"0": false,
				"1": false,
				"2": true,
				"3": false,
				"4": false
			},
			"answer_index": 2
		}`

		w := PerformRequest(r, "POST", "/api/v1/decks/"+deckID+"/cards/", strings.NewReader(body), token2)
		if w.Code != 201 {
			t.Errorf("Expected status code 201, got %d", w.Code)
		}

		resp := w.Body.String()
		expectedSubstring := `"question":"What is 2 + 2?"`
		if !strings.Contains(resp, expectedSubstring) {
			t.Errorf("Expected response body to contain %q, got %q", expectedSubstring, resp)
		}
	})

	// Delete one card from the deck
	t.Run("Delete one card from the deck", func(t *testing.T) {
		// First, get the list of cards to find a card ID to delete
		w := PerformRequest(r, "GET", "/api/v1/decks/"+deckID, nil, token1)
		if w.Code != 200 {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}

		var respData struct {
			Cards []struct {
				ID string `json:"id"`
			} `json:"cards"`
		}

		if err := json.Unmarshal(w.Body.Bytes(), &respData); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if len(respData.Cards) == 0 {
			t.Fatal("No cards found in deck, cannot perform delete test")
		}

		cardID := respData.Cards[0].ID

		// Now we can delete the card
		w = PerformRequest(r, "DELETE", "/api/v1/decks/"+deckID+"/cards/"+cardID, nil, token1)
		if w.Code != 204 {
			t.Errorf("Expected status code 204, got %d", w.Code)
		}
	})

	t.Run("Delete the created user with the deck and cards", func(t *testing.T) {
		w := PerformRequest(r, "DELETE", "/api/v1/users/", nil, token1)
		if w.Code != 204 {
			t.Errorf("Expected status code 204, got %d", w.Code)
		}
	})

	// Shared deck should now be empty
	t.Run("Verify shared decks are empty after owner deletion", func(t *testing.T) {
		w := PerformRequest(r, "GET", "/api/v1/users/decks/shared", nil, token2)
		if w.Code != 200 {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}
		resp := w.Body.String()
		expected := `[]`
		if resp != expected {
			t.Errorf("Expected response body %q, got %q", expected, resp)
		}
	})
}
