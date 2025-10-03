package firebase

import (
	"context"
	"memora/internal/config"
	"memora/internal/errors"
	"memora/internal/utils"

	"cloud.google.com/go/firestore"
)

type CardRepository interface {
	CreateCard(ctx context.Context, card any) (string, error)
	GetCard(ctx context.Context, id string) (map[string]any, error)
	UpdateCard(ctx context.Context, firestoreUpdates []firestore.Update, id string) error
	DeleteCard(ctx context.Context, id string) error
}

type FirestoreCardRepo struct {
	client *firestore.Client
}

func NewFirestoreCardRepo(client *firestore.Client) *FirestoreCardRepo {
	return &FirestoreCardRepo{client: client}
}

func (r *FirestoreCardRepo) GetCard(ctx context.Context, id string) (map[string]any, error) {
	doc, err := r.client.Collection(config.CardsCollection).Doc(id).Get(ctx)
	if err != nil {
		return nil, errors.ErrInvalidId
	}

	return doc.Data(), nil
}

func (r *FirestoreCardRepo) CreateCard(ctx context.Context, card any) (string, error) {
	return utils.AddToDB(r.client, ctx, config.CardsCollection, card)
}

func (r *FirestoreCardRepo) UpdateCard(
	ctx context.Context,
	firestoreUpdates []firestore.Update,
	id string,
) error {
	return utils.UpdateDocumentInDB(
		r.client,
		ctx,
		config.CardsCollection,
		id,
		firestoreUpdates,
	)
}

func (r *FirestoreCardRepo) DeleteCard(ctx context.Context, id string) error {
	return utils.DeleteDocumentInDB(
		r.client,
		ctx,
		config.CardsCollection,
		id,
	)
}
