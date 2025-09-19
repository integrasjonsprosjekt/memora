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
	UpdateDeck(ctx context.Context, id string, update models.UpdateDeck) error
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

	returnID, _, err := r.client.Collection(internal.DECKS_COLLECTION).Add(ctx, deck)
	return returnID.ID, err
}

func (r *FirestoreDeckRepo) GetOneDeck(ctx context.Context, id string) (models.DeckResponse, error) {
	var deck models.Deck
	var response models.DeckResponse
	var cards []models.Card

	doc, err := r.client.Collection(internal.DECKS_COLLECTION).Doc(id).Get(ctx)
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

func (r *FirestoreDeckRepo) UpdateDeck(ctx context.Context, id string, update models.UpdateDeck) error {
	deckRef := r.client.Collection(internal.DECKS_COLLECTION).Doc(id)
	deckSnap, err := deckRef.Get(ctx)
	if err != nil {
		return customerror.ErrNotFound
	}

	var deck models.Deck
	if err := deckSnap.DataTo(&deck); err != nil {
		return err
	}

	for _, c := range update.Cards {
		if c.ID == "" {
			cardRef := r.client.Collection(internal.CARDS_COLLECTION).NewDoc()
			_, err := cardRef.Set(ctx, models.Card{
				Type:    c.Type,
				Front:   c.Front,
				Back:    c.Back,
				Options: c.Options,
				Answer:  c.Answer,
			})
			if err != nil {
				return err
			}
			deck.Cards = append(deck.Cards, cardRef)
		} else {
			var cardRef *firestore.DocumentRef
			for _, ref := range deck.Cards {
				if ref.ID == c.ID {
					cardRef = ref
					break
				}
			}

			if cardRef == nil {
				return customerror.ErrNotFound
			}

			if c.Type == "" && c.Front == "" && c.Back == "" && len(c.Options) == 0 && c.Answer == "" {
				newCards := []*firestore.DocumentRef{}
				for _, ref := range deck.Cards {
					if ref.ID != c.ID {
						newCards = append(newCards, ref)
					}
				}
				deck.Cards = newCards
				_, _ = cardRef.Delete(ctx)
			} else {
				updates := map[string]interface{}{}

				if c.Type != "" {
					updates["type"] = c.Type
				}
				if c.Front != "" {
					updates["front"] = c.Front
				}
				if c.Back != "" {
					updates["back"] = c.Back
				}
				if c.Options != nil {
					updates["options"] = c.Options
				}
				if c.Answer != "" {
					updates["answer"] = c.Answer
				}

				_, err := cardRef.Set(ctx, updates, firestore.MergeAll)
				if err != nil {
					return err
				}
			}
		}
	}
	_, err = deckRef.Set(ctx, map[string]interface{}{"cards": deck.Cards}, firestore.MergeAll)
	return err
}
