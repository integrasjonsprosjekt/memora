package firebase

import "cloud.google.com/go/firestore"

// Repositories groups all Firestore repositories.
type Repositories struct {
	User *FirestoreUserRepo
	Card *FirestoreCardRepo
	Deck *FirestoreDeckRepo
}

// NewRepositories creates a new Repositories struct with the provided Firestore client.
func NewRepositories(client *firestore.Client) *Repositories {
	return &Repositories{
		User: NewFirestoreUserRepo(client),
		Card: NewFirestoreCardRepo(client),
		Deck: NewFirestoreDeckRepo(client),
	}
}
