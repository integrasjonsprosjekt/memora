package firebase

import (
	"context"
	"memora/internal/config"

	"cloud.google.com/go/firestore"
)

type CardRepository interface {
	CreateCard(ctx context.Context, card any) (string, error)
}

type FirestoreCardRepo struct {
	client *firestore.Client
}

func NewFirestoreCardRepo(client *firestore.Client) *FirestoreCardRepo {
	return &FirestoreCardRepo{client: client}
}

func (r *FirestoreCardRepo) CreateCard(ctx context.Context, card any) (string, error) {
	id, _, err := r.client.Collection(config.CardsCollection).Add(ctx, card)
	return id.ID, err
}
