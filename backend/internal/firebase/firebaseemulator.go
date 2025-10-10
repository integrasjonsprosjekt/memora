// internal/firebase/init_emulator.go
package firebase

import (
	"context"
	"fmt"

	"memora/internal/config"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

func InitEmulator() (*firestore.Client, error) {
	config.GetEnv("FIRESTORE_EMULATOR_HOST", config.GetEnv("FIRESTORE_EMULATOR_HOST", "localhost:8081"))

	ctx := context.Background()

	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing firestore client: %v", err)
	}

	return client, nil
}
