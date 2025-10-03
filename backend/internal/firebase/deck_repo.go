package firebase

import (
	"context"
	"memora/internal/config"
	customerror "memora/internal/customError"
	"memora/internal/models"

	"cloud.google.com/go/firestore"
)

type FirestoreDeckRepo struct {
	client *firestore.Client
}

type DeckRepository interface {
	AddDeck(ctx context.Context, deck models.CreateDeck) (string, error)
	GetOneDeck(ctx context.Context, id string) (models.DeckResponse, error)
	UpdateDeck(ctx context.Context, id string, update models.UpdateDeck) error
}

func NewFirestoreDeckRepo(client *firestore.Client) *FirestoreDeckRepo {
	return &FirestoreDeckRepo{client: client}
}

// AddDeck implements DeckRepository.
func (r *FirestoreDeckRepo) AddDeck(ctx context.Context, deck models.CreateDeck) (string, error) {
	_, err := r.client.Collection(config.UsersCollection).Doc(deck.OwnerID).Get(ctx)
	if err != nil {
		return "", customerror.ErrNotFound
	}

	returnID, _, err := r.client.Collection(config.DecksCollection).Add(ctx, deck)
	return returnID.ID, err
}

func (r *FirestoreDeckRepo) GetOneDeck(
	ctx context.Context,
	id string,
) (models.DeckResponse, error) {
	return models.DeckResponse{}, nil
}

func (r *FirestoreDeckRepo) UpdateDeck(
	ctx context.Context,
	id string,
	update models.UpdateDeck,
) error {
	return nil
}
