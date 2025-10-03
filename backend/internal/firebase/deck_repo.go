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
}

func NewFirestoreDeckRepo(client *firestore.Client) *FirestoreDeckRepo {
	return &FirestoreDeckRepo{client: client}
}

// AddDeck implements DeckRepository.
func (r *FirestoreDeckRepo) AddDeck(ctx context.Context, deck models.CreateDeck) (string, error) {
	_, err := r.client.Collection(config.UsersCollection).Doc(deck.OwnerID).Get(ctx)
	if err != nil {
		return "", errors.ErrNotFound
	}

	returnID, _, err := r.client.Collection(config.DecksCollection).Add(ctx, deck)
	return returnID.ID, err
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
	deckExists, err := utils.CheckIfDocumentExists(r.client, ctx, config.DecksCollection, deckID)
	if err != nil {
		return err
	}

	cardExists, err := utils.CheckIfDocumentExists(r.client, ctx, config.CardsCollection, cardID)
	if err != nil {
		return err
	}

	if !deckExists || !cardExists {
		return errors.ErrInvalidId
	}

	deckRef := r.client.Collection(config.DecksCollection).Doc(deckID)
	cardRef := r.client.Collection(config.CardsCollection).Doc(cardID)

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
	deckExists, err := utils.CheckIfDocumentExists(r.client, ctx, config.DecksCollection, deckID)
	if err != nil {
		return err
	}

	cardExists, err := utils.CheckIfDocumentExists(r.client, ctx, config.CardsCollection, cardID)
	if err != nil {
		return err
	}

	if !deckExists || !cardExists {
		return errors.ErrInvalidId
	}

	deckRef := r.client.Collection(config.DecksCollection).Doc(deckID)
	cardRef := r.client.Collection(config.CardsCollection).Doc(cardID)

	_, err = deckRef.Update(ctx, []firestore.Update{
		{Path: "cards", Value: firestore.ArrayUnion(cardRef)},
	})
	if err != nil {
		return err
	}

	return nil
}
