package firebase

import (
	"context"
	"memora/internal/config"
	"memora/internal/errors"

	"cloud.google.com/go/firestore"
)

type CardRepository interface {
	CreateCard(ctx context.Context, card any) (string, error)
	GetCard(ctx context.Context, id string) (map[string]any, error)
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

func (r *FirestoreCardRepo) GetCard(ctx context.Context, id string) (map[string]any, error) {
	doc, err := r.client.Collection(config.CardsCollection).Doc(id).Get(ctx)
	if err != nil {
		return nil, errors.ErrInvalidId
	}

	return doc.Data(), nil
}
