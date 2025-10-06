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
	RemoveCardsFromDeck(ctx context.Context, deckID string, cardIDs []string) error
	AddCardsToDeck(ctx context.Context, deckID string, cardIDs []string) error
	UpdateDeck(ctx context.Context, firestoreUpdates []firestore.Update, id string) error
	RemoveEmailsFromShared(ctx context.Context, deckID string, emails []string) error
	AddEmailsToShared(ctx context.Context, deckID string, emails []string) error
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

func (r *FirestoreDeckRepo) RemoveCardsFromDeck(
	ctx context.Context,
	deckID string,
	cardIDs []string,
) error {
	deckSnap, err := utils.GetDocumentIfExists(r.client, ctx, config.DecksCollection, deckID)
	if err != nil {
		return err
	}
	deckRef := deckSnap.Ref
	return r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		cardsIface := make([]any, len(cardIDs))
		for i, id := range cardIDs {
			cardSnap, err := utils.GetDocumentIfExists(
				r.client,
				ctx,
				config.CardsCollection,
				id,
			)
			if err != nil {
				return errors.ErrFailedUpdatingCards
			}
			cardsIface[i] = cardSnap.Ref
		}
		return tx.Update(deckRef, []firestore.Update{
			{Path: "cards", Value: firestore.ArrayRemove(cardsIface...)},
		})
	})
}

func (r *FirestoreDeckRepo) AddCardsToDeck(
	ctx context.Context,
	deckID string,
	cardIDs []string,
) error {
	deckSnap, err := utils.GetDocumentIfExists(r.client, ctx, config.DecksCollection, deckID)
	if err != nil {
		return err
	}
	deckRef := deckSnap.Ref
	return r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		cardsIface := make([]any, len(cardIDs))
		for i, id := range cardIDs {
			cardSnap, err := utils.GetDocumentIfExists(
				r.client,
				ctx,
				config.CardsCollection,
				id,
			)
			if err != nil {
				return errors.ErrFailedUpdatingCards
			}
			cardsIface[i] = cardSnap.Ref
		}
		return tx.Update(deckRef, []firestore.Update{
			{Path: "cards", Value: firestore.ArrayUnion(cardsIface...)},
		})
	})
}

func (r *FirestoreDeckRepo) AddEmailsToShared(
	ctx context.Context,
	deckID string,
	emails []string,
) error {
	deckRef := r.client.Collection(config.DecksCollection).Doc(deckID)
	return r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		for _, email := range emails {
			exists, err := utils.UserExistsByEmail(r.client, ctx, email)
			if err != nil {
				return errors.ErrFailedUpdatingEmail
			}
			if !exists {
				return errors.ErrInvalidEmailNotPresent
			}
		}

		emailsIface := make([]any, len(emails))
		for i, v := range emails {
			emailsIface[i] = v
		}
		return tx.Update(deckRef, []firestore.Update{
			{Path: "shared_emails", Value: firestore.ArrayUnion(emailsIface...)},
		})
	})
}

func (r *FirestoreDeckRepo) RemoveEmailsFromShared(
	ctx context.Context,
	deckID string,
	emails []string,
) error {
	deckRef := r.client.Collection(config.DecksCollection).Doc(deckID)
	return r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		for _, email := range emails {
			exists, err := utils.UserExistsByEmail(r.client, ctx, email)
			if err != nil {
				return errors.ErrFailedUpdatingEmail
			}
			if !exists {
				return errors.ErrInvalidEmailNotPresent
			}
		}

		emailsIface := make([]any, len(emails))
		for i, v := range emails {
			emailsIface[i] = v
		}
		return tx.Update(deckRef, []firestore.Update{
			{Path: "shared_emails", Value: firestore.ArrayRemove(emailsIface...)},
		})
	})
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
