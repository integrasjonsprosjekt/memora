package firebase

import (
	"context"
	"memora/internal/config"
	"memora/internal/errors"
	"memora/internal/models"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type UserRepository interface {
	AddUser(ctx context.Context, u models.CreateUser) (string, error)
	GetUser(ctx context.Context, id string) (models.User, error)
}

type FirestoreUserRepo struct {
	client *firestore.Client
}

func NewFirestoreUserRepo(client *firestore.Client) *FirestoreUserRepo {
	return &FirestoreUserRepo{client: client}
}

func (r *FirestoreUserRepo) GetUser(ctx context.Context, id string) (models.User, error) {
	var userStruct = models.User{}
	user, err := r.client.Collection(config.UsersCollection).Doc(id).Get(ctx)
	if err != nil {
		return models.User{}, errors.ErrNotFound
	}
	if err := user.DataTo(&userStruct); err != nil {
		return models.User{}, err
	}

	userStruct.ID = id

	return userStruct, nil
}

func (r *FirestoreUserRepo) AddUser(ctx context.Context, u models.CreateUser) (string, error) {
	iter := r.client.Collection(config.UsersCollection).Where("email", "==", u.Email).Limit(1).Documents(ctx)
	doc, err := iter.Next()
	if err != nil && err != iterator.Done {
		return "", nil
	}
	if doc != nil {
		return "", errors.ErrInvalidUser
	}

	id, _, err := r.client.Collection(config.UsersCollection).Add(ctx, u)
	return id.ID, err
}
