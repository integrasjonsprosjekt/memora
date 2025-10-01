package firebase

import (
	"cloud.google.com/go/firestore"
)

type CardRepository interface {
}

type FirestoreCardRepo struct {
	client *firestore.Client
}

func NewFirestoreCardRepo(client *firestore.Client) *FirestoreCardRepo {
	return &FirestoreCardRepo{client: client}
}
