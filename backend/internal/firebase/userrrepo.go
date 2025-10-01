package firebase

import "cloud.google.com/go/firestore"

type UserRepository interface {
}

type FirestoreUserRepo struct {
	client *firestore.Client
}

func NewFirestoreUserRepo(client *firestore.Client) *FirestoreUserRepo {
	return &FirestoreUserRepo{client: client}
}
