package utils

import (
	"context"
	"memora/internal/config"
	"memora/internal/errors"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func FetchByID[T any](
	client *firestore.Client,
	ctx context.Context,
	collection, id string,
) (T, error) {
	var responseStruct T
	response, err := client.Collection(collection).Doc(id).Get(ctx)
	if err != nil {
		return responseStruct, errors.ErrNotFound
	}
	if err := response.DataTo(&responseStruct); err != nil {
		return responseStruct, err
	}
	return responseStruct, nil
}

func AddToDB[T any](
	client *firestore.Client,
	ctx context.Context,
	collection string,
	data T,
) (string, error) {
	id, _, err := client.Collection(collection).Add(ctx, data)
	return id.ID, err
}

func UpdateDocumentInDB(
	client *firestore.Client,
	ctx context.Context,
	collection, id string,
	updates []firestore.Update,
) error {
	docRef := client.Collection(collection).Doc(id)
	_, err := docRef.Get(ctx)
	if err != nil {
		return errors.ErrInvalidId
	}

	_, err = docRef.Update(ctx, updates)
	if err != nil {
		return err
	}
	return nil
}

func DeleteDocumentInDB(
	client *firestore.Client,
	ctx context.Context,
	collection, id string,
) error {
	docRef := client.Collection(collection).Doc(id)
	_, err := docRef.Get(ctx)
	if err != nil {
		return errors.ErrInvalidId
	}

	_, err = docRef.Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

func GetDocumentIfExists(
	client *firestore.Client,
	ctx context.Context,
	collection, id string,
) (*firestore.DocumentSnapshot, error) {
	snap, err := client.Collection(collection).Doc(id).Get(ctx)
	if err != nil {
		return nil, errors.ErrInvalidId
	}
	if !snap.Exists() {
		return nil, errors.ErrNotFound
	}
	return snap, nil
}

func UserExistsByEmail(
	client *firestore.Client,
	ctx context.Context,
	email string,
) (bool, error) {
	iter := client.Collection(config.UsersCollection).
		Where("email", "==", email).
		Limit(1).
		Documents(ctx)
	doc, err := iter.Next()
	if err != nil && err != iterator.Done {
		return false, err
	}
	return doc != nil, nil
}
