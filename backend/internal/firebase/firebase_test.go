package firebase

import (
	"os"
	"testing"
)

// setTestCredentials sets the GOOGLE_APPLICATION_CREDENTIALS environment variable
// and returns a cleanup function that restores the original value
// If credPath is nil, the environment variable is unset
func setTestCredentials(credPath *string) func() {
	originalValue, exists := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")

	if credPath == nil {
		os.Unsetenv("GOOGLE_APPLICATION_CREDENTIALS")
	} else {
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", *credPath)
	}

	return func() {
		if exists {
			os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", originalValue)
		} else {
			os.Unsetenv("GOOGLE_APPLICATION_CREDENTIALS")
		}
	}
}

func TestInit(t *testing.T) {
	t.Run("missing GOOGLE_APPLICATION_CREDENTIALS", func(t *testing.T) {
		defer setTestCredentials(nil)()

		client, err := Init()

		if err == nil {
			t.Error("Expected error when GOOGLE_APPLICATION_CREDENTIALS is not set, got nil")
		}

		expectedErrorMsg := "GOOGLE_APPLICATION_CREDENTIALS not set"
		if err.Error() != expectedErrorMsg {
			t.Errorf("Expected error message '%s', got '%s'", expectedErrorMsg, err.Error())
		}

		if client != nil {
			t.Error("Expected nil client when GOOGLE_APPLICATION_CREDENTIALS is not set")
		}
	})

	t.Run("with GOOGLE_APPLICATION_CREDENTIALS set", func(t *testing.T) {
		// Set a test value (this will still fail Firebase initialization, but passes the env var check)
		testCredPath := "/tmp/test-credentials.json"
		defer setTestCredentials(&testCredPath)()

		client, err := Init()

		// We expect this to fail because the credentials file doesn't exist,
		// but it should NOT fail with our specific "GOOGLE_APPLICATION_CREDENTIALS not set" error
		if err != nil && err.Error() == "GOOGLE_APPLICATION_CREDENTIALS not set" {
			t.Error("Should not fail with 'GOOGLE_APPLICATION_CREDENTIALS not set' error when env var is set")
		}

		// Since we're using a fake credentials file, the client should be nil due to Firebase initialization failure
		if client != nil {
			// If we somehow get a client, make sure to close it
			client.Close()
		}
	})
}

// TestInitWithValidCredentials is a helper test that can be run when you have valid Firebase credentials
// To run this test, set the GOOGLE_APPLICATION_CREDENTIALS environment variable to the path of valid credentials
func TestInitWithValidCredentials(t *testing.T) {
	credPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credPath == "" {
		t.Skip("Skipping test with valid credentials. Set GOOGLE_APPLICATION_CREDENTIALS to run this test.")
	}

	defer setTestCredentials(&credPath)()

	client, err := Init()

	if err != nil {
		t.Fatalf("Expected successful initialization with valid credentials, got error: %v", err)
	}

	if client == nil {
		t.Fatal("Expected non-nil client with valid credentials")
	}

	defer client.Close()
}

// BenchmarkInit benchmarks the Init function
func BenchmarkInit(t *testing.B) {
	// This benchmark will only run meaningful tests if credentials are available
	credPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credPath == "" {
		t.Skip("Skipping benchmark. Set GOOGLE_APPLICATION_CREDENTIALS to run this benchmark.")
	}

	defer setTestCredentials(&credPath)()

	t.ResetTimer()

	for i := 0; i < t.N; i++ {
		client, err := Init()
		if err != nil {
			t.Fatalf("Benchmark failed: %v", err)
		}
		if client != nil {
			client.Close()
		}
	}
}
