package firebase

import (
	"context"
	"memora/internal/config"
	"memora/internal/errors"
	"memora/internal/models"
	"memora/internal/utils"

	"cloud.google.com/go/firestore"
)

type FirestoreDeckRepo struct {
	client *firestore.Client
}

type DeckRepository interface {
	AddDeck(ctx context.Context, deck models.CreateDeck) (string, error)
	GetOneDeck(ctx context.Context, id string) (models.Deck, error)
	RemoveCardFromDeck(ctx context.Context, deckID, cardID string) error
	AddCardToDeck(ctx context.Context, deckID, cardID string) error
	UpdateDeck(ctx context.Context, firestoreUpdates []firestore.Update, id string) error
	RemoveEmailFromShared(ctx context.Context, email, deckID string) error
	AddEmailToShared(ctx context.Context, email, deckID string) error
	DeleteDeck(ctx context.Context, id string) error
}

func NewFirestoreDeckRepo(client *firestore.Client) *FirestoreDeckRepo {
	return &FirestoreDeckRepo{client: client}
}

func (r *FirestoreDeckRepo) AddDeck(ctx context.Context, deck models.CreateDeck) (string, error) {
	_, err := utils.GetDocumentIfExists(r.client, ctx, config.UsersCollection, deck.OwnerID)
	if err != nil {
		return "", err
	}

	for _, email := range deck.SharedEmails {
		exists, err := utils.UserExistsByEmail(r.client, ctx, email)
		if err != nil {
			return "", err
		}
		if !exists {
			return "", errors.ErrInvalidEmailNotPresent
		}
	}

	return utils.AddToDB(r.client, ctx, config.DecksCollection, deck)
}

func (r *FirestoreDeckRepo) GetOneDeck(
	ctx context.Context,
	id string,
) (models.Deck, error) {
	return utils.FetchByID[models.Deck](r.client, ctx, config.DecksCollection, id)
}

func (r *FirestoreDeckRepo) RemoveCardFromDeck(
	ctx context.Context,
	deckID, cardID string,
) error {
	deckSnap, err := utils.GetDocumentIfExists(r.client, ctx, config.DecksCollection, deckID)
	if err != nil {
		return err
	}

	cardSnap, err := utils.GetDocumentIfExists(r.client, ctx, config.CardsCollection, cardID)
	if err != nil {
		return err
	}

	deckRef := deckSnap.Ref
	cardRef := cardSnap.Ref
	_, err = deckRef.Update(ctx, []firestore.Update{
		{Path: "cards", Value: firestore.ArrayRemove(cardRef)},
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *FirestoreDeckRepo) AddCardToDeck(
	ctx context.Context,
	deckID, cardID string,
) error {
	deckSnap, err := utils.GetDocumentIfExists(r.client, ctx, config.DecksCollection, deckID)
	if err != nil {
		return err
	}

	cardSnap, err := utils.GetDocumentIfExists(r.client, ctx, config.CardsCollection, cardID)
	if err != nil {
		return err
	}

	deckRef := deckSnap.Ref
	cardRef := cardSnap.Ref
	_, err = deckRef.Update(ctx, []firestore.Update{
		{Path: "cards", Value: firestore.ArrayUnion(cardRef)},
	})

	return err
}

func (r *FirestoreDeckRepo) AddEmailToShared(
	ctx context.Context,
	email, deckID string,
) error {
	exists, err := utils.UserExistsByEmail(r.client, ctx, email)
	if err != nil {
		return err
	}
	if !exists {
		return errors.ErrInvalidEmailNotPresent
	}

	deckRef := r.client.Collection(config.DecksCollection).Doc(deckID)

	_, err = deckRef.Update(ctx, []firestore.Update{
		{Path: "shared_emails", Value: firestore.ArrayUnion(email)},
	})
	return err
}

func (r *FirestoreDeckRepo) RemoveEmailFromShared(
	ctx context.Context,
	email, deckID string,
) error {
	exists, err := utils.UserExistsByEmail(r.client, ctx, email)
	if err != nil {
		return err
	}
	if !exists {
		return errors.ErrInvalidEmailNotPresent
	}

	deckRef := r.client.Collection(config.DecksCollection).Doc(deckID)

	_, err = deckRef.Update(ctx, []firestore.Update{
		{Path: "shared_emails", Value: firestore.ArrayRemove(email)},
	})
	return err
}

func (r *FirestoreDeckRepo) UpdateDeck(
	ctx context.Context,
	firestoreUpdates []firestore.Update,
	id string,
) error {
	return utils.UpdateDocumentInDB(
		r.client,
		ctx,
		config.DecksCollection,
		id,
		firestoreUpdates,
	)
}

func (r *FirestoreDeckRepo) DeleteDeck(
	ctx context.Context,
	id string,
) error {
	return utils.DeleteDocumentInDB(
		r.client,
		ctx,
		config.DecksCollection,
		id,
	)
}