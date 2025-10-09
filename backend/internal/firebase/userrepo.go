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
	AddUser(ctx context.Context, u models.CreateUser) (string, error)

	// GetUser fetches a user from Firestore by ID.
	// Error on failure or if the ID is invalid.
	// Returns the user on success.
	GetUser(ctx context.Context, id string) (models.User, error)

	// GetDecksOwned fetches all decks owned by a user.
	// Error on failure or if the user ID is invalid.
	// Returns the decks ID and title on success.
	GetDecksOwned(ctx context.Context, id string) ([]models.DisplayDeck, error)

	// GetDecksShared fetches all decks shared with a user.
	// Error on failure or if the user ID is invalid.
	// Returns the decks ID and title on success.
	GetDecksShared(ctx context.Context, id string) ([]models.DisplayDeck, error)

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
) (models.User, error) {
	user, err := utils.FetchByID[models.User](r.client, ctx, config.UsersCollection, id)
	if err != nil {
		return user, err
	}

	user.ID = id

	return user, nil
}

// GetUser fetches a user from Firestore by ID.
// Error on failure or if the ID is invalid.
// Returns the user on success.
func (r *FirestoreUserRepo) GetDecksOwned(
	ctx context.Context,
	id string,
) ([]models.DisplayDeck, error) {
	var decks []models.DisplayDeck

	_, err := utils.GetDocumentIfExists(r.client, ctx, config.UsersCollection, id)
	if err != nil {
		return nil, err
	}

	// Get all decks owned by the user.
	iter := r.client.Collection(config.DecksCollection).
		Where("owner_id", "==", id).
		Documents(ctx)

	// Append each deck to the slice.
	for {
		doc, err := iter.Next()

		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var deck models.DisplayDeck
		// Map the document data to the DisplayDeck struct.
		if err := doc.DataTo(&deck); err != nil {
			return nil, err
		}

		// Set the ID and append to the slice.
		deck.ID = doc.Ref.ID
		decks = append(decks, deck)
	}

	return decks, nil
}

// GetDecksOwned fetches all decks owned by a user.
// Error on failure or if the user ID is invalid.
// Returns the decks ID and title on success.
func (r *FirestoreUserRepo) GetDecksShared(
	ctx context.Context,
	id string,
) ([]models.DisplayDeck, error) {
	var decks []models.DisplayDeck

	// Get the user by ID.
	user, err := utils.FetchByID[models.User](r.client, ctx, config.UsersCollection, id)
	if err != nil {
		return nil, err
	}

	// Create iterator where shared_emails array contains the user's email.
	iter := r.client.Collection(config.DecksCollection).
		Where("shared_emails", "array-contains", user.Email).
		Documents(ctx)

	// Loop through the documents and append to the decks slice.
	for {
		doc, err := iter.Next()

		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var deck models.DisplayDeck
		// Map the document data to the DisplayDeck struct.
		if err := doc.DataTo(&deck); err != nil {
			return nil, err
		}

		// Set the ID and append to the slice.
		deck.ID = doc.Ref.ID
		decks = append(decks, deck)
	}

	return decks, nil
}

// AddUser adds a new user to Firestore.
// Error on failure or if the email is already present.
// Returns the new user's ID on success.
func (r *FirestoreUserRepo) AddUser(
	ctx context.Context, 
	user models.CreateUser,
) (string, error) {
	// Check if the email is already present.
	exists, err := utils.UserExistsByEmail(r.client, ctx, user.Email)
	if err != nil {
		return "", err
	}
	if exists {
		return "", errors.ErrInvalidEmailPresent
	}

	// Email is unique, add the user to Firestore.
	return utils.AddToDB(r.client, ctx, config.UsersCollection, user)
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
func (r *FirestoreUserRepo) DeleteUser(
	ctx context.Context, 
	id string,
) error {
	return utils.DeleteDocumentInDB(
		r.client,
		ctx,
		config.UsersCollection,
		id,
	)
}
