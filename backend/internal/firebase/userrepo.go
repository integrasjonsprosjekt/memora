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

type UserRepository interface {
	AddUser(ctx context.Context, u models.CreateUser) (string, error)
	GetUser(ctx context.Context, id string) (models.User, error)
	UpdateUser(ctx context.Context, firestoreUpdates []firestore.Update, id string) error
	DeleteUser(ctx context.Context, id string) error
}

type FirestoreUserRepo struct {
	client *firestore.Client
}

func NewFirestoreUserRepo(client *firestore.Client) *FirestoreUserRepo {
	return &FirestoreUserRepo{client: client}
}

func (r *FirestoreUserRepo) GetUser(ctx context.Context, id string) (models.User, error) {
	user, err := utils.FetchByID[models.User](r.client, ctx, config.UsersCollection, id)
	if err != nil {
		return user, err
	}

	user.ID = id

	return user, nil
}

func (r *FirestoreUserRepo) AddUser(ctx context.Context, u models.CreateUser) (string, error) {
	iter := r.client.Collection(config.UsersCollection).
		Where("email", "==", u.Email).
		Limit(1).
		Documents(ctx)
	doc, err := iter.Next()
	if err != nil && err != iterator.Done {
		return "", err
	}
	if doc != nil {
		return "", errors.ErrInvalidUser
	}

	return utils.AddToDB(r.client, ctx, config.UsersCollection, u)
}

func (r *FirestoreUserRepo) UpdateUser(
	ctx context.Context,
	firestoreUpdates []firestore.Update,
	id string,
) error {
	docRef := r.client.Collection(config.UsersCollection).Doc(id)

	_, err := docRef.Get(ctx)
	if err != nil {
		return errors.ErrInvalidId
	}

	_, err = docRef.Update(ctx, firestoreUpdates)
	if err != nil {
		return err
	}
	return nil
}

func (r *FirestoreUserRepo) DeleteUser(ctx context.Context, id string) error {
	docRef := r.client.Collection(config.UsersCollection).Doc(id)
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
