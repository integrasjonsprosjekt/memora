package firebase

import (
	"cloud.google.com/go/firestore"
	"github.com/redis/go-redis/v9"
)

// Repositories groups all Firestore repositories.
type Repositories struct {
	User  *FirestoreUserRepo
	Card  *FirestoreCardRepo
	Deck  *FirestoreDeckRepo
	Auth  *FirebaseAuthRepo
	Redis *redis.Client
}

// NewRepositories creates a new Repositories struct with the provided Firestore client.
func NewRepositories(client *firestore.Client, auth *FirebaseAuthRepo, redis *redis.Client) *Repositories {
	return &Repositories{
		User:  NewFirestoreUserRepo(client),
		Card:  NewFirestoreCardRepo(client),
		Deck:  NewFirestoreDeckRepo(client),
		Auth:  auth,
		Redis: redis,
	}
}
