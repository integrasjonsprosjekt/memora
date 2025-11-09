package firebase

import (
	"context"
	"fmt"
	"memora/internal/config"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/redis/go-redis/v9"
)

// Init initializes and returns a Firestore client.
// It requires the GOOGLE_APPLICATION_CREDENTIALS environment variable to be set.
// Error if initialization fails.
func Init() (*firestore.Client, *firebase.App, *redis.Client, error) {
	if _, ok := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS"); !ok {
		return nil, nil, nil, fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS not set")
	}

	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error initializing app: %v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error initializing firestore client: %v", err)
	}

	rbd := redis.NewClient(&redis.Options{
		Addr:     config.GetEnv("REDIS_ADDR", "localhost:6379"),
		Password: config.GetEnv("REDIS_PASSWORD", ""),
		DB:       0,
	})

	if _, err := rbd.Ping(ctx).Result(); err != nil {
		return nil, nil, nil, fmt.Errorf("error connecting to redis: %v", err)
	}

	return client, app, rbd, nil
}
