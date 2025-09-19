package firebase

import (
	"context"
	"memora/internal"
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
}

func NewFirestoreDeckRepo(client *firestore.Client) *FirestoreDeckRepo {
	return &FirestoreDeckRepo{client: client}
}

// AddDeck implements DeckRepository.
func (r *FirestoreDeckRepo) AddDeck(ctx context.Context, deck models.CreateDeck) (string, error) {
	_, err := r.client.Collection(internal.USERS_COLLECTION).Doc(deck.OwnerID).Get(ctx)
	if err != nil {
		return "", customerror.ErrNotFound
	}

	returnID, _, err := r.client.Collection(internal.DECK_COLLECTION).Add(ctx, deck)
	return returnID.ID, err
}

func (r *FirestoreDeckRepo) GetOneDeck(ctx context.Context, id string) (models.DeckResponse, error) {
	var deck models.Deck
	var response models.DeckResponse
	var cards []models.Card

	doc, err := r.client.Collection(internal.DECK_COLLECTION).Doc(id).Get(ctx)
	if err != nil {
		return response, customerror.ErrNotFound
	}

	if err := doc.DataTo(&deck); err != nil {
		return response, err
	}

	for _, ref := range deck.Cards {
		snap, err := ref.Get(ctx)
		if err != nil {
			return response, err
		}

		var c models.Card
		if err := snap.DataTo(&c); err != nil {
			return response, err
		}

		cards = append(cards, models.Card{
			ID:      snap.Ref.ID,
			Type:    c.Type,
			Front:   c.Front,
			Back:    c.Back,
			Options: c.Options,
			Answer:  c.Answer,
		})
	}

	return models.DeckResponse{
		ID:           doc.Ref.ID,
		OwnerID:      deck.OwnerID,
		SharedEmails: deck.SharedEmails,
		Cards:        cards,
	}, nil
}
