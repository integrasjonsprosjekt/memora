package firebase

import (
	"context"
	"memora/internal"
	"memora/internal/models"

	"cloud.google.com/go/firestore"
)

type FirestoreDeckRepo struct {
	client *firestore.Client
}

type DeckRepository interface {
	AddDeck(ctx context.Context, deck models.CreateDeck) (string, error)
}

func NewFirestoreDeckRepo(client *firestore.Client) *FirestoreDeckRepo {
	return &FirestoreDeckRepo{client: client}
}

// AddDeck implements DeckRepository.
func (r *FirestoreDeckRepo) AddDeck(ctx context.Context, deck models.CreateDeck) (string, error) {
	_, err := r.client.Collection(internal.USERS_COLLECTION).Doc(deck.Owner_id).Get(ctx)
	if err != nil {
		return "", models.ErrUserNotFound
	}

	returnID, _, err := r.client.Collection(internal.DECK_COLLECTION).Add(ctx, deck)
	return returnID.ID, err
}
