package firebase

import (
	"context"

	"cloud.google.com/go/firestore"
)

type FirestoreDeckRepo struct {
	client *firestore.Client
}

type DeckRepository interface {
	AddDeck(ctx context.Context) (string, error)
}

func NewFirestoreDeckRepo(client *firestore.Client) *FirestoreDeckRepo {
	return &FirestoreDeckRepo{client: client}
}

// AddDeck implements DeckRepository.
func (f *FirestoreDeckRepo) AddDeck(ctx context.Context) (string, error) {
	panic("unimplemented")
}