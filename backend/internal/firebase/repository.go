package firebase

import "cloud.google.com/go/firestore"

type Repositories struct {
	User *FirestoreUserRepo
	Card *FirestoreCardRepo
}

func NewRepositories(client *firestore.Client) *Repositories {
	return &Repositories{
		User: NewFirestoreUserRepo(client),
		Card: NewFirestoreCardRepo(client),
	}
}
