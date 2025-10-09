package firebase

import (
	"context"
	"memora/internal/config"
	"memora/internal/errors"
	"memora/internal/models"
	"memora/internal/utils"

	"cloud.google.com/go/firestore"
)

// DeckRepository methods used for storing, updating and deleting data
type DeckRepository interface {
	// AddDeck adds a new deck to firestore.
	// Error on failure or malformed data, returns ID on success
	AddDeck(ctx context.Context, deck models.CreateDeck) (string, error)

	// GetOneDeck fetches an existing deck from firestore.
	// Error on fail or if the ID is invalid, returns deck on success
	GetOneDeck(ctx context.Context, id string) (models.Deck, error)

	// RemoveCardsFromDeck removes given cards from a deck based on their respected IDs.
	// Error on failure in transaction, nil on success
	RemoveCardsFromDeck(ctx context.Context, deckID string, cardIDs []string) error

	// AddCardsToDeck adds given cards into a deck based on their respected IDs.
	// Error on failure in transaction, nil on success
	AddCardsToDeck(ctx context.Context, deckID string, cardIDs []string) error

	// UpdateDeck updates everything except emails and cards in a given deck.
	// Error on failure, or if ID is invalid, nil on success
	UpdateDeck(ctx context.Context, firestoreUpdates []firestore.Update, id string) error

	// RemoveEmailsFromShared removes given emails from the decks shared emails.
	// Error on failure in transaction, nil on success
	RemoveEmailsFromShared(ctx context.Context, deckID string, emails []string) error

	// AddEmailsToShared adds given email into the decks shared emails.
	// Error on failure in transaction, nil on success
	AddEmailsToShared(ctx context.Context, deckID string, emails []string) error

	// DeleteDeck deletes a given deck from firestore.
	// Error on failure, or if ID is invalid, nil on success
	DeleteDeck(ctx context.Context, id string) error
}

// FirestoreDeckRepo holds the database connection needed for fetching
type FirestoreDeckRepo struct {
	client *firestore.Client
}

// NewFirestoreDeckRepo creates a new repoistory and returns pointer to struct
func NewFirestoreDeckRepo(client *firestore.Client) *FirestoreDeckRepo {
	return &FirestoreDeckRepo{client: client}
}

// AddDeck checks if the owned ID is valid, and then adds the decks data into firestore.
// Error on failure, or if parameters SharedEmails and OwnerID is invalid.
// Returns decks ID on success
func (r *FirestoreDeckRepo) AddDeck(ctx context.Context, deck models.CreateDeck) (string, error) {

	// Check if the user exists
	_, err := utils.GetDocumentIfExists(r.client, ctx, config.UsersCollection, deck.OwnerID)
	if err != nil {
		return "", err
	}

	// Loop over email and check if they exist.
	// If it fails on one email it returns error
	for _, email := range deck.SharedEmails {
		exists, err := utils.UserExistsByEmail(r.client, ctx, email)
		if err != nil {
			return "", err
		}
		if !exists {
			return "", errors.ErrInvalidEmailNotPresent
		}
	}

	// Safely add the deck to firestore
	return utils.AddToDB(r.client, ctx, config.DecksCollection, deck)
}

// GetOneDeck from firestore.
// Error on failure, or if deck is invalid.
// Returns deck on success
func (r *FirestoreDeckRepo) GetOneDeck(
	ctx context.Context,
	id string,
) (models.Deck, error) {
	return utils.FetchByID[models.Deck](r.client, ctx, config.DecksCollection, id)
}

// RemoveCardsFromDeck removes any given cards from a deck by IDs
// Error on any failure in transaction
// Returns nil on success
func (r *FirestoreDeckRepo) RemoveCardsFromDeck(
	ctx context.Context,
	deckID string,
	cardIDs []string,
) error {
	// Make sure the deck exists
	deckSnap, err := utils.GetDocumentIfExists(r.client, ctx, config.DecksCollection, deckID)
	if err != nil {
		return err
	}
	deckRef := deckSnap.Ref

	// Define transaction when updating the cards
	return r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// Interface to hold the cards refrence to be removed
		cardsIface := make([]any, len(cardIDs))

		for i, id := range cardIDs {
			// Check if the card actually exists
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
		// Remove it from the array
		return tx.Update(deckRef, []firestore.Update{
			{Path: "cards", Value: firestore.ArrayRemove(cardsIface...)},
		})
	})
}

// AddCardsToDeck adds existing cards to a deck by their IDs.
// Error on failure at any point when its processing the cards IDs.
// Returns nil on success.
func (r *FirestoreDeckRepo) AddCardsToDeck(
	ctx context.Context,
	deckID string,
	cardIDs []string,
) error {
	// Make sure the deck exists
	deckSnap, err := utils.GetDocumentIfExists(r.client, ctx, config.DecksCollection, deckID)
	if err != nil {
		return err
	}
	deckRef := deckSnap.Ref

	// Define transaction to add cards to the collection
	return r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {

		// Interface used to update
		cardsIface := make([]any, len(cardIDs))
		for i, id := range cardIDs {
			// Make sure the card exists
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

		// Update the cards in the deck
		return tx.Update(deckRef, []firestore.Update{
			{Path: "cards", Value: firestore.ArrayUnion(cardsIface...)},
		})
	})
}

// AddEmailsToShared adds emails to a deck to gain permissions on the deck.
// Error on failure, or if email does not exist.
// Returns nil on success
func (r *FirestoreDeckRepo) AddEmailsToShared(
	ctx context.Context,
	deckID string,
	emails []string,
) error {
	// Check if the deck exists
	deckRef := r.client.Collection(config.DecksCollection).Doc(deckID)

	// Run update transactions
	return r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		for _, email := range emails {
			exists, err := utils.UserExistsByEmail(r.client, ctx, email)
			// Unforseen error
			if err != nil {
				return errors.ErrFailedUpdatingEmail
			}
			// Does not exist
			if !exists {
				return errors.ErrInvalidEmailNotPresent
			}
		}

		emailsIface := make([]any, len(emails))
		for i, v := range emails {
			emailsIface[i] = v
		}

		// Update the shared emails in firestore
		return tx.Update(deckRef, []firestore.Update{
			{Path: "shared_emails", Value: firestore.ArrayUnion(emailsIface...)},
		})
	})
}

// RemoveEmailsFromShared removes emails from a decks shared emails.
// Error on failure, or if email does not exist.
// Returns nil on success
func (r *FirestoreDeckRepo) RemoveEmailsFromShared(
	ctx context.Context,
	deckID string,
	emails []string,
) error {
	// Check if the deck exists
	deckRef := r.client.Collection(config.DecksCollection).Doc(deckID)

	// Run update transaction
	return r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		for _, email := range emails {
			exists, err := utils.UserExistsByEmail(r.client, ctx, email)
			// Unforseen error
			if err != nil {
				return errors.ErrFailedUpdatingEmail
			}
			// Does not exist
			if !exists {
				return errors.ErrInvalidEmailNotPresent
			}
		}

		emailsIface := make([]any, len(emails))
		for i, v := range emails {
			emailsIface[i] = v
		}

		// Update the shared emails in firestore
		return tx.Update(deckRef, []firestore.Update{
			{Path: "shared_emails", Value: firestore.ArrayRemove(emailsIface...)},
		})
	})
}

// DeleteDeck deletes a deck from firestore based on its ID.
// Error on failure, or if ID is invalid.
// Returns nil on success
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

// DeleteDeck deletes a deck from firestore based on its ID.
// Error on failure, or if ID is invalid.
// Returns nil on success
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
