package firebase

import (
	"context"
	"memora/internal/config"
	"memora/internal/errors"
	"memora/internal/models"

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

	// GetCardProgress retrieves the progress of a card for a specific user.
	// Error on fail, returns the progress if successful
	GetCardProgress(ctx context.Context, deckID, cardID, userID string) (models.CardProgress, error)

	UpdateProgress(
		ctx context.Context,
		deckID, cardID, userID string,
		firestoreUpdates models.CardProgress,
	) error

	GetDueCardsInDeck(
		ctx context.Context,
		deckID, userID string,
		limit int,
		cursor string,
	) ([]map[string]any, string, bool, error)
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

// GetCardProgress retrieves the progress of a card for a specific user.
// Returns the progress or an error if the operation fails.
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

// UpdateProgress updates the progress of a card for a specific user.
// Returns an error if the operation fails.
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

// GetDueCardsInDeck fetches due cards for a user in a deck with pagination support
// Returns a list of cards, next cursor, hasMore flag, and an error if the operation fails.
func (r *FirestoreCardRepo) GetDueCardsInDeck(
	ctx context.Context,
	deckID, userID string,
	limit int,
	cursor string,
) ([]map[string]any, string, bool, error) {

	var cards []map[string]any
	var nextCursor string

	// Fetch all progress documents for the user to identify studied cards
	progressMap := make(map[string]bool)
	allProgressIter := r.client.
		Collection(config.DecksCollection).Doc(deckID).
		Collection(config.UsersCollection).Doc(userID).
		Collection(config.ProgressCollection).
		Documents(ctx)

	for {
		doc, err := allProgressIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, "", false, err
		}
		progressMap[doc.Ref.ID] = true
	}

	allProgressIter.Stop()

	// Prioritize unstudied cards first
	unstudiedQuery := r.client.
		Collection(config.DecksCollection).Doc(deckID).
		Collection(config.CardsCollection).
		OrderBy(firestore.DocumentID, firestore.Asc).
		Limit(limit + 1)

	// Apply a cursor if provided
	if cursor != "" && cursor[:10] == "unstudied_" {
		unstudiedQuery = unstudiedQuery.StartAfter(cursor[10:])
	}

	unstudiedIter := unstudiedQuery.Documents(ctx)
	defer unstudiedIter.Stop()

	lastUnstudiedID := ""
	skipped := 0

	// Iterate through unstudied cards
	for {
		cardDoc, err := unstudiedIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, "", false, err
		}

		cardID := cardDoc.Ref.ID
		lastUnstudiedID = cardID

		// If the card is unstudied, add it to the results
		if !progressMap[cardID] {
			cardData := cardDoc.Data()
			cardData["id"] = cardID
			cards = append(cards, cardData)
		} else {
			// Skip studied cards used to decide if it has more unstudied
			skipped++
		}
	}

	// Check if we have enough unstudied cards
	hasMoreUnstudied := len(cards)+skipped > limit
	if hasMoreUnstudied {
		// Get the new pagination cursor and return
		nextCursor = "unstudied_" + lastUnstudiedID
		return cards[:limit], nextCursor, true, nil
	}

	// We need more cards, fetch studied cards based on due date
	if len(cards) < limit {
		remaining := limit - len(cards)

		// Build the due date query
		dueQuery := r.client.
			Collection(config.DecksCollection).Doc(deckID).
			Collection(config.UsersCollection).Doc(userID).
			Collection(config.ProgressCollection).
			OrderBy("due", firestore.Asc).
			Limit(remaining + 1)

		// Apply cursor if provided in previous query
		if len(cursor) >= 4 && cursor[:4] == "due_" {
			dueQuery = dueQuery.StartAfter(cursor[4:])
		}

		dueIter := dueQuery.Documents(ctx)

		defer dueIter.Stop()

		lastDueID := ""

		cardRefs := []*firestore.DocumentRef{}

		// Batch fetch due cards to not load network multiple times
		for {
			progressDoc, err := dueIter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, "", false, err
			}

			cardID := progressDoc.Ref.ID
			lastDueID = cardID

			cardRef := r.client.
				Collection(config.DecksCollection).Doc(deckID).
				Collection(config.CardsCollection).Doc(cardID)

			cardRefs = append(cardRefs, cardRef)
		}

		// If there are card references, fetch their data
		if len(cardRefs) > 0 {
			cardDocs, err := r.client.GetAll(ctx, cardRefs)
			if err != nil {
				return nil, "", false, err
			}

			for _, cardDoc := range cardDocs {
				cardData := cardDoc.Data()
				cardData["id"] = cardDoc.Ref.ID
				cards = append(cards, cardData)
			}
		}

		// Check if we have more studied cards for pagination
		if len(cards) > limit {
			nextCursor = "due_" + lastDueID
			return cards[:limit], nextCursor, true, nil
		}
	}

	// No more cards to paginate through
	return cards, "", false, nil
}
