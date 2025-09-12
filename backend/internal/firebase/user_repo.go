package firebase

import (
	"context"
	"errors"
	"fmt"
	"memora/internal"
	"memora/internal/models"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type FirestoreUserRepo struct {
	client *firestore.Client
}

var (
	ErrUserNotFound = errors.New("user not found")
)

// GetUser implements UserRepository.
func (r *FirestoreUserRepo) GetUser(ctx context.Context, id string) (models.User, error) {
	var userStruct = models.User{}
	user, err := r.client.Collection(internal.USERS_COLLECTION).Doc(id).Get(ctx)
	if err != nil {
		return models.User{}, ErrUserNotFound
	}
	if err := user.DataTo(&userStruct); err != nil {
		return models.User{}, fmt.Errorf("unable to marshal user")
	}
	return userStruct, nil
}

func NewFirestoreUserRepo(client *firestore.Client) *FirestoreUserRepo {
	return &FirestoreUserRepo{client: client}
}

func (r *FirestoreUserRepo) AddUser(ctx context.Context, u models.CreateUser) (string, error) {
	iter := r.client.Collection(internal.USERS_COLLECTION).Where("email", "==", u.Email).Limit(1).Documents(ctx)
	doc, err := iter.Next()
	if err != nil && err != iterator.Done {
		return "", nil
	}
	if doc != nil {
		return "", fmt.Errorf("user with this email exists alredy")
	}

	id, _, err := r.client.Collection(internal.USERS_COLLECTION).Add(ctx, u)
	return id.ID, err
}
