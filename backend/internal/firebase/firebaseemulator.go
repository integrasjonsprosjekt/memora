// internal/firebase/init_emulator.go
package firebase

import (
	"context"
	"os"

	"memora/internal/config"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

func InitEmulator() (*firestore.Client, error) {
	config.Init() // Ensure config is initialized

	host := config.GetEnv("FIRESTORE_EMULATOR_HOST", "localhost:8081")

	err := os.Setenv("FIRESTORE_EMULATOR_HOST", host)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	// Project ID is required but can be any string
	client, err := firestore.NewClient(ctx, "demo-project", option.WithoutAuthentication())
	if err != nil {
		return nil, err
	}

	return client, nil
}
