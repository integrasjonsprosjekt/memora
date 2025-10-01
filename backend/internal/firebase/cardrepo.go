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
	UpdateCard(ctx context.Context, firestoreUpdates []firestore.Update, id string) error
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
	id, _, err := r.client.Collection(config.CardsCollection).Add(ctx, card)
	return id.ID, err
}

func (r *FirestoreCardRepo) UpdateCard(ctx context.Context, firestoreUpdates []firestore.Update, id string) error {
	docRef := r.client.Collection(config.CardsCollection).Doc(id)

	_, err := docRef.Get(ctx)
	if err != nil {
		return errors.ErrInvalidId
	}

	_, err = docRef.Update(ctx, firestoreUpdates)
	if err != nil {
		return err
	}
	return nil
}