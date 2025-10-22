package firebase

import (
	"context"
	"memora/internal/config"
	"memora/internal/errors"
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
	// Error on failure or if the user ID is invalid.
	// Returns the decks ID and title on success.
	GetDecks(ctx context.Context, id string, fields []string) ([]models.DisplayDeck, error)

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
) ([]models.DisplayDeck, error) {

	// Get the user by ID. After middleware is introduced, this can be omitted.
	user, err := utils.FetchByID[models.User](r.client, ctx, config.UsersCollection, id, fields)
	if err != nil {
		return nil, err
	}

	// Get all decks owned by the user.
	iter := r.client.Collection(config.DecksCollection).
		Where("owner_id", "==", id).
		Documents(ctx)
	decksOwned, err := utils.ReadDataFromIterator[models.DisplayDeck](iter)
	if err != nil {
		return nil, err
	}

	// Create iterator where shared_emails array contains the user's email.
	iter = r.client.Collection(config.DecksCollection).
		Where("shared_emails", "array-contains", user.Email).
		Documents(ctx)
	decksShared, err := utils.ReadDataFromIterator[models.DisplayDeck](iter)
	if err != nil {
		return nil, err
	}

	decks := append(decksOwned, decksShared...)

	return decks, nil
}

// AddUser adds a new user to Firestore.
// Error on failure or if the email is already present.
// Returns the new user's ID on success.
func (r *FirestoreUserRepo) AddUser(
	ctx context.Context,
	user models.CreateUser,
	id string,
) error {
	// Check if the email is already present.
	exists, err := utils.UserExistsByEmail(r.client, ctx, user.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.ErrInvalidEmailPresent
	}

	// Email is unique, add the user to Firestore.
	_, err = r.client.Collection(config.UsersCollection).Doc(id).Set(ctx, user)
	return err
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
