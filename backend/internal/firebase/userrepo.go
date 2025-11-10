package firebase

import (
	"context"
	"memora/internal/config"
	"memora/internal/models"
	"memora/internal/utils"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// UserRepository defines the interface for user-related Firestore operations.
type UserRepository interface {
	// AddUser adds a new user to Firestore.
	// Error on failure or if the email is already present.
	// Returns the new user's ID on success.
	AddUser(ctx context.Context, u models.CreateUser, id string) error

	// GetUser fetches a user from Firestore by ID.
	// Error on failure or if the ID is invalid.
	// Returns the user on success.
	GetUser(ctx context.Context, id string, fields []string) (models.User, error)

	// GetDecks fetches all decks for a user.
	// GetDecks fetches all decks for a user.
	// Error on failure or if the user ID is invalid.
	// Returns the decks ID and title on success.
	GetDecks(ctx context.Context, id string, fields []string) (models.UserDecks, error)

	// UpdateUser updates fields of an existing user in Firestore.
	// Error on failure or if the ID is invalid.
	// Returns nil on success.
	UpdateUser(ctx context.Context, firestoreUpdates []firestore.Update, id string) error

	// DeleteUser deletes a user from Firestore by ID.
	// Error on failure or if the ID is invalid.
	// Returns nil on success.
	DeleteUser(ctx context.Context, id string) error
}

// FirestoreUserRepo implements the UserRepository interface using Firestore as the backend.
type FirestoreUserRepo struct {
	client *firestore.Client
}

// NewFirestoreUserRepo creates and returns a pointer to the FirestoreUserRepo.
func NewFirestoreUserRepo(client *firestore.Client) *FirestoreUserRepo {
	return &FirestoreUserRepo{client: client}
}

// AddUser adds a new user to Firestore.
// Error on failure or if the email is already present.
// Returns the new user's ID on success.
func (r *FirestoreUserRepo) GetUser(
	ctx context.Context,
	id string,
	fields []string,
) (models.User, error) {
	user, err := utils.FetchByID[models.User](r.client, ctx, config.UsersCollection, id, fields)
	if err != nil {
		return user, err
	}

	user.ID = id

	return user, nil
}

// GetUser fetches a user from Firestore by ID.
// Error on failure or if the ID is invalid.
// Returns the user on success.
func (r *FirestoreUserRepo) GetDecks(
	ctx context.Context,
	id string,
	fields []string,
) (models.UserDecks, error) {

	// Get the user by ID. After middleware is introduced, this can be omitted.
	user, err := utils.FetchByID[models.User](
		r.client,
		ctx,
		config.UsersCollection,
		id,
		[]string{"email"},
	)
	if err != nil {
		return models.UserDecks{}, err
	}

	// Channel to receive results from goroutines
	type result struct {
		decks []models.DisplayDeck
		err   error
	}
	ownedChan := make(chan result, 1)
	sharedChan := make(chan result, 1)

	go func() {
		// Get all decks owned by the user.
		iter := r.client.Collection(config.DecksCollection).
			Where("owner_id", "==", id).
			Select(fields...).
			Documents(ctx)
		decksOwned, err := readDataFromIterator(iter)
		ownedChan <- result{decks: decksOwned, err: err}
	}()

	go func() {
		// Create iterator where shared_emails array contains the user's email.
		iter := r.client.Collection(config.DecksCollection).
			Where("shared_emails", "array-contains", user.Email).
			Select(fields...).
			Documents(ctx)
		decksShared, err := readDataFromIterator(iter)
		sharedChan <- result{decks: decksShared, err: err}
	}()

	ownedRes := <-ownedChan
	if ownedRes.err != nil {
		return models.UserDecks{}, ownedRes.err
	}

	sharedRes := <-sharedChan
	if sharedRes.err != nil {
		return models.UserDecks{}, sharedRes.err
	}

	decks := models.UserDecks{
		OwnedDecks:  ownedRes.decks,
		SharedDecks: sharedRes.decks,
	}

	return decks, nil
}

// AddUser adds a new user to Firestore, with a mock deck.
// Error on failure or if the email is already present.
// Returns the new user's ID on success.
func (r *FirestoreUserRepo) AddUser(
	ctx context.Context,
	user models.CreateUser,
	id string,
) error {
	return r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		if err := tx.Set(r.client.Collection(config.UsersCollection).Doc(id), user); err != nil {
			return err
		}

		deckRef := r.client.Collection(config.DecksCollection).NewDoc()
		mockDeck := models.CreateDeck{
			Title:        "Default Deck",
			OwnerID:      id,
			SharedEmails: []string{},
		}
		if err := tx.Set(deckRef, mockDeck); err != nil {
			return err
		}

		cardsCollection := deckRef.Collection(config.CardsCollection)
		mockCards := getMockCards()
		for _, card := range mockCards {
			cardRef := cardsCollection.NewDoc()
			if err := tx.Set(cardRef, card); err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateUser updates fields of an existing user in Firestore.
// Error on failure or if the ID is invalid.
// Returns nil on success.
func (r *FirestoreUserRepo) UpdateUser(
	ctx context.Context,
	firestoreUpdates []firestore.Update,
	id string,
) error {
	return utils.UpdateDocumentInDB(
		r.client,
		ctx,
		config.UsersCollection,
		id,
		firestoreUpdates,
	)
}

// DeleteUser deletes a user from Firestore by ID.
// Error on failure or if the ID is invalid.
// Returns nil on success
// TODO refactor this to support batching
func (r *FirestoreUserRepo) DeleteUser(
	ctx context.Context,
	id string,
) error {
	decksIter := r.client.Collection(config.DecksCollection).
		Where("owner_id", "==", id).
		Documents(ctx)
	defer decksIter.Stop()

	bulkWriter := r.client.BulkWriter(ctx)

	for {
		deckDoc, err := decksIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		// Delete all cards in this deck
		cardsIter := deckDoc.Ref.Collection("cards").Documents(ctx)

		for {
			cardDoc, err := cardsIter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}

			// Schedule delete using BulkWriter
			if _, err := bulkWriter.Delete(cardDoc.Ref); err != nil {
				return err
			}
		}

		// Delete the deck itself
		if _, err := bulkWriter.Delete(deckDoc.Ref); err != nil {
			return err
		}
	}

	userRef := r.client.Collection(config.UsersCollection).Doc(id)

	if _, err := bulkWriter.Delete(userRef); err != nil {
		return err
	}

	// Wait for all operations to complete
	bulkWriter.Flush()
	bulkWriter.End()

	return nil
}

// readDataFromIterator reads documents from a Firestore DocumentIterator
// and converts them into a slice of DisplayDeck models.
// Used to read decks owned or shared with a user
func readDataFromIterator(iter *firestore.DocumentIterator) ([]models.DisplayDeck, error) {
	var results []models.DisplayDeck

	// Iterate through the documents and convert them to DisplayDeck models
	for {
		doc, err := iter.Next()
		if err != nil {
			// Break the loop if there are no more documents
			if err == iterator.Done {
				break
			}
			return nil, err
		}

		// Convert document data to DisplayDeck model
		var item models.DisplayDeck
		if err := doc.DataTo(&item); err != nil {
			return nil, err
		}

		// Set the document ID and append to results
		item.ID = doc.Ref.ID
		results = append(results, item)
	}

	return results, nil
}

// getMockCards returns the default cards for a new user's deck
func getMockCards() []any {
	return []any{
		models.FrontBackCard{
			Front: "Welcome to Memora!",
			Back:  "This is your first flashcard. Edit or delete it to get started.",
			Type:  utils.FRONT_BACK_CARD,
		},
		models.MultipleChoiceCard{
			Question: "What is Memora?",
			Options: map[string]bool{
				"A flashcard app":         true,
				"A social media platform": false,
				"A video game":            false,
			},
			Type: utils.MULTIPLE_CHOICE_CARD,
		},
		models.OrderedCard{
			Question: "Arrange the steps to create a deck in order.",
			Options:  []string{"Create an account", "Add a deck", "Add cards to the deck"},
			Type:     utils.ORDERED_CARD,
		},
		models.BlanksCard{
			Question: "Memora is a {} app for {}.",
			Answers:  []string{"flashcard", "learning"},
			Type:     utils.BLANKS_CARD,
		},
	}
}
