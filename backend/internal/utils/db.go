package utils

import (
	"context"
	"memora/internal/config"
	"memora/internal/errors"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// FetchByID retrieves a document by its ID from the specified collection
// and maps it to the provided generic type T. Returns an error if the document
// is not found or if there is an issue during retrieval or mapping.
func FetchByID[T any](
	client *firestore.Client,
	ctx context.Context,
	collection, id string,
) (T, error) {
	// Initialize a variable of type T to hold the response
	var responseStruct T

	// Retrieve the document from Firestore
	response, err := client.Collection(collection).Doc(id).Get(ctx)

	// Handle errors during retrieval
	if err != nil {
		return responseStruct, errors.ErrNotFound
	}

	// Map the document data to the response struct
	if err := response.DataTo(&responseStruct); err != nil {
		return responseStruct, err
	}

	return responseStruct, nil
}

// AddToDB adds a new document to the specified collection in Firestore.
// It takes a generic type T for the document data and returns the new
// document's ID or an error if the operation fails.
func AddToDB[T any](
	client *firestore.Client,
	ctx context.Context,
	collection string,
	data T,
) (string, error) {
	id, _, err := client.Collection(collection).Add(ctx, data)
	return id.ID, err
}

// UpdateDocumentInDB updates an existing document in the specified collection
// in Firestore. It takes a list of Firestore updates and returns an error
// if the operation fails or if the document does not exist.
func UpdateDocumentInDB(
	client *firestore.Client,
	ctx context.Context,
	collection, id string,
	updates []firestore.Update,
) error {
	// Check if the document exists
	docRef := client.Collection(collection).Doc(id)
	_, err := docRef.Get(ctx)
	if err != nil {
		return errors.ErrInvalidId
	}

	// Perform the update
	_, err = docRef.Update(ctx, updates)
	if err != nil {
		return err
	}
	return nil
}

// DeleteDocumentInDB deletes a document from the specified collection in Firestore.
// It returns an error if the operation fails or if the document does not exist.
func DeleteDocumentInDB(
	client *firestore.Client,
	ctx context.Context,
	collection, id string,
) error {
	// Check if the document exists
	docRef := client.Collection(collection).Doc(id)
	_, err := docRef.Get(ctx)
	if err != nil {
		return errors.ErrInvalidId
	}

	// Perform the deletion
	_, err = docRef.Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

// GetDocumentIfExists retrieves a document from Firestore by its ID
// and checks if it exists. Returns the document snapshot or an error
// if the document does not exist or if there is an issue during retrieval.
func GetDocumentIfExists(
	client *firestore.Client,
	ctx context.Context,
	collection, id string,
) (*firestore.DocumentSnapshot, error) {
	// Retrieve the document snapshot
	snap, err := client.Collection(collection).Doc(id).Get(ctx)
	if err != nil {
		return nil, errors.ErrInvalidId
	}

	// Check if the document exists
	if !snap.Exists() {
		return nil, errors.ErrNotFound
	}
	return snap, nil
}

// UserExistsByEmail checks if a user with the specified email exists in Firestore.
// Returns true if the user exists, false otherwise, along with any error encountered.
func UserExistsByEmail(
	client *firestore.Client,
	ctx context.Context,
	email string,
) (bool, error) {
	// Query Firestore for a user with the given email
	iter := client.Collection(config.UsersCollection).
		Where("email", "==", email).
		Limit(1).
		Documents(ctx)
	
	// Check if any document was returned
	doc, err := iter.Next()
	if err != nil && err != iterator.Done {
		return false, err
	}
	return doc != nil, nil
}
