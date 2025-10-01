package firebase

import "cloud.google.com/go/firestore"

type Repositories struct {
	User *FirestoreUserRepo
}

func NewRepositories(client *firestore.Client) *Repositories {
	return &Repositories{
		User: NewFirestoreUserRepo(client),
	}
}
