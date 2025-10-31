package firebase

import (
	"context"
	"memora/internal/config"
	"memora/internal/errors"
	"memora/internal/models"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// CardStore defines methods for storing, retrieving and updating cards.
type CardRepository interface {
	// CreateCard adds a new card into firestore.
	// Error on fail, returns the ID if succesfull
	CreateCard(ctx context.Context, card any, deckID string) (string, error)

	// GetCardsInDeck fetches all cards in a given deck with cursor-based pagination.
	// cursor is the ID of the last card from the previous page (empty string for first page)
	// Error on fail, returns a list of cards on success
	GetCardsInDeck(
		ctx context.Context,
		deckID string,
		limit int,
		cursor string,
	) ([]map[string]any, bool, error)

	// GetCard returns the raw data for a card for a given ID.
	// Error on fail, or if ID is not valid
	GetCardInDeck(ctx context.Context, deckID, cardID string) (map[string]any, error)

	// UpdateCard updates an existing card in firestore.
	// Error on fail or if the ID is not valid, nil on success
	UpdateCard(
		ctx context.Context,
		firestoreUpdates []firestore.Update,
		deckID, cardID string,
	) error

	// DeleteCard deletes an existing card in firestore.
	// Error on fail or if the ID is not valid, nil on success
	DeleteCard(ctx context.Context, deckID, cardID string) error

	// CreateProgress creates a new progress entry for a card and user.
	// Error on fail, returns the ID if successful
	CreateProgress(ctx context.Context, deckID, cardID, userID string, progress models.CardProgress) (string, error)

	// GetCardProgress retrieves the progress of a card for a specific user.
	// Error on fail, returns the progress if successful
	GetCardProgress(ctx context.Context, deckID, cardID, userID string) (models.CardProgress, error)

	UpdateProgress(ctx context.Context, deckID, cardID, userID string, firestoreUpdates models.CardProgress) error

	GetDueCardsInDeck(ctx context.Context, deckID, userID string, limit int) ([]map[string]any, error)
}

// FirestoreCardRepo holds the connection to the database
type FirestoreCardRepo struct {
	client *firestore.Client
}

// NewFirestoreCardRepo creates and returns a pointer to the repository
func NewFirestoreCardRepo(client *firestore.Client) *FirestoreCardRepo {
	return &FirestoreCardRepo{client: client}
}

// GetCardsInDeck fetches all cards in a given deck with cursor-based pagination.
// cursor is the document ID of the last card from the previous page (empty for first page)
// Returns a list of cards and an error if the operation fails.
func (r *FirestoreCardRepo) GetCardsInDeck(
	ctx context.Context,
	deckID string,
	limit int,
	cursor string,
) ([]map[string]any, bool, error) {
	// Build the base query
	query := r.client.Collection(config.DecksCollection).
		Doc(deckID).
		Collection(config.CardsCollection).
		OrderBy(firestore.DocumentID, firestore.Asc).
		Limit(limit + 1) // Fetch one extra to check for more pages

	// If we have a cursor, start after that document ID
	// Using just the ID value is more efficient than fetching the document
	if cursor != "" {
		query = query.StartAfter(cursor)
	}

	iter := query.Documents(ctx)
	defer iter.Stop()

	var result []map[string]any

	// Iterate through the documents in the collection
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, false, err
		}

		// Append the document data along with its ID
		data := doc.Data()
		data["id"] = doc.Ref.ID

		result = append(result, data)
	}

	hasMore := false
	if len(result) > limit {
		hasMore = true
		result = result[:limit] // Trim the extra document used for pagination check
	}

	return result, hasMore, nil
}

// GetCard takes a context and id, and returns the raw data for a card
// error if the card can not be fetched
func (r *FirestoreCardRepo) GetCardInDeck(
	ctx context.Context,
	deckID, cardID string,
) (map[string]any, error) {
	doc, err := r.client.Collection(config.DecksCollection).
		Doc(deckID).
		Collection(config.CardsCollection).
		Doc(cardID).
		Get(ctx)
	if err != nil {
		return nil, errors.ErrInvalidId
	}

	return doc.Data(), nil
}

// CreateCard takes a context and a card, adds it to the database, and
// returns the created card or an error if the operation fails.
func (r *FirestoreCardRepo) CreateCard(
	ctx context.Context,
	card any, deckID string,
) (string, error) {
	docRef, _, err := r.client.Collection(config.DecksCollection).
		Doc(deckID).
		Collection(config.CardsCollection).
		Add(ctx, card)
	if err != nil {
		return "", err
	}
	return docRef.ID, nil
}

// UpdateCard takes a context, an update payload, and an ID, and updates
// the corresponding card in the database. It returns an error if the update
// fails or the card cannot be found
func (r *FirestoreCardRepo) UpdateCard(
	ctx context.Context,
	firestoreUpdates []firestore.Update,
	deckID, cardID string,
) error {
	docRef := r.client.
		Collection(config.DecksCollection).
		Doc(deckID).
		Collection(config.CardsCollection).
		Doc(cardID)

	// Perform the update
	_, err := docRef.Update(ctx, firestoreUpdates)
	if err != nil {
		return err
	}
	return nil
}

// DeleteCard takes a context and ID, and deletes the corresponding
// card in the database. It returns an error if the delete fails,
// or if it cannot be found
func (r *FirestoreCardRepo) DeleteCard(
	ctx context.Context,
	deckID, cardID string,
) error {
	docRef := r.client.
		Collection(config.DecksCollection).
		Doc(deckID).
		Collection(config.CardsCollection).
		Doc(cardID)
	_, err := docRef.Get(ctx)
	if err != nil {
		return errors.ErrInvalidId
	}

	// Delete the document
	_, err = docRef.Delete(ctx)
	return err
}

func (r *FirestoreCardRepo) CreateProgress(
	ctx context.Context,
	deckID, cardID, userID string,
	progress models.CardProgress,
) (string, error) {
	_, err := r.client.
		Collection(config.DecksCollection).Doc(deckID).
		Collection(config.UsersCollection).Doc(userID).
		Collection(config.ProgressCollection).Doc(cardID).
		Set(ctx, progress)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (r *FirestoreCardRepo) GetCardProgress(
	ctx context.Context,
	deckID, cardID, userID string,
) (models.CardProgress, error) {
	doc, err := r.client.
		Collection(config.DecksCollection).Doc(deckID).
		Collection(config.UsersCollection).Doc(userID).
		Collection(config.ProgressCollection).Doc(cardID).
		Get(ctx)
	if err != nil {
		return models.CardProgress{}, errors.ErrInvalidId
	}

	var progress models.CardProgress
	if err := doc.DataTo(&progress); err != nil {
		return models.CardProgress{}, err
	}

	return progress, nil
}

func (r *FirestoreCardRepo) UpdateProgress(
	ctx context.Context,
	deckID, cardID, userID string,
	firestoreUpdates models.CardProgress,
) error {
	docRef := r.client.
		Collection(config.DecksCollection).Doc(deckID).
		Collection(config.UsersCollection).Doc(userID).
		Collection(config.ProgressCollection).Doc(cardID)

	_, err := docRef.Set(ctx, firestoreUpdates)
	if err != nil {
		return err
	}
	return nil
}

func (r *FirestoreCardRepo) GetDueCardsInDeck(
	ctx context.Context,
	deckID, userID string,
	limit int,
) ([]map[string]any, error) {
	now := time.Now()

	progressQuery := r.client.
		Collection(config.DecksCollection).Doc(deckID).
		Collection(config.UsersCollection).Doc(userID).
		Collection(config.ProgressCollection).
		Where("due", "<=", now).
		OrderBy("due", firestore.Asc).
		Limit(limit)

	iter := progressQuery.Documents(ctx)
	defer iter.Stop()

	var dueCards []map[string]any

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		cardID := doc.Ref.ID

		cardDoc, err := r.client.
			Collection(config.DecksCollection).Doc(deckID).
			Collection(config.CardsCollection).Doc(cardID).
			Get(ctx)
		if err != nil {
			return nil, err
		}

		cardData := cardDoc.Data()
		cardData["id"] = cardDoc.Ref.ID
		dueCards = append(dueCards, cardData)
	}

	if len(dueCards) < limit {
		remaining := limit - len(dueCards)

		allCardsIter := r.client.
			Collection(config.DecksCollection).Doc(deckID).
			Collection(config.CardsCollection).
			Limit(remaining).
			Documents(ctx)

		for {
			if len(dueCards) >= limit {
				allCardsIter.Stop()
				break
			}

			cardDoc, err := allCardsIter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, err
			}

			cardID := cardDoc.Ref.ID

			_, err = r.client.
				Collection(config.DecksCollection).Doc(deckID).
				Collection(config.UsersCollection).Doc(userID).
				Collection(config.ProgressCollection).Doc(cardID).
				Get(ctx)

			if err != nil {
				cardData := cardDoc.Data()
				cardData["id"] = cardDoc.Ref.ID
				dueCards = append(dueCards, cardData)
			}
		}
		allCardsIter.Stop()
	}
	return dueCards, nil
}
