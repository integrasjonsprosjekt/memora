// internal/firebase/init_emulator.go
package firebase

import (
	"context"
	"os"

	"memora/internal/config"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func InitEmulator() (*firestore.Client, *firebase.App, error) {
	config.Init() // Ensure config is initialized

	host := config.GetEnv("FIRESTORE_EMULATOR_HOST", "127.0.0.1:8081")
	projectID := config.GetEnv("FIRESTORE_PROJECT_ID", "memora-test")

	err := os.Setenv("FIRESTORE_EMULATOR_HOST", host)
	if err != nil {
		return nil, nil, err
	}

	err = os.Setenv("FIREBASE_AUTH_EMULATOR_HOST", "127.0.0.1:9099")
	if err != nil {
		return nil, nil, err
	}

	conf := &firebase.Config{
		ProjectID: projectID,
	}

	app, err := firebase.NewApp(context.Background(), conf, option.WithoutAuthentication())
	if err != nil {
		return nil, nil, err
	}

	// Project ID is required but can be any string
	client, err := app.Firestore(context.Background())
	if err != nil {
		return nil, nil, err
	}

	return client, app, nil
}
