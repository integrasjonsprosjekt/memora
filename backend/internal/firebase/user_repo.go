package firebase

import (
	"context"
	"fmt"
	"memora/internal"
	"memora/internal/models"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type FirestoreUserRepo struct {
	client *firestore.Client
}

type UserRepository interface {
	AddUser(ctx context.Context, u models.CreateUser) (string, error)
	GetUser(ctx context.Context, id string) (models.User, error)
	UpdateUser(ctx context.Context, update map[string]interface{}, id string) error
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
		return "", models.ErrInvalidUser
	}

	id, _, err := r.client.Collection(internal.USERS_COLLECTION).Add(ctx, u)
	return id.ID, err
}

func (r *FirestoreUserRepo) GetUser(ctx context.Context, id string) (models.User, error) {
	var userStruct = models.User{}
	user, err := r.client.Collection(internal.USERS_COLLECTION).Doc(id).Get(ctx)
	if err != nil {
		return models.User{}, models.ErrUserNotFound
	}
	if err := user.DataTo(&userStruct); err != nil {
		return models.User{}, fmt.Errorf("unable to marshal user")
	}

	userStruct.ID = id

	return userStruct, nil
}

func (r *FirestoreUserRepo) UpdateUser(ctx context.Context, u map[string]interface{}, id string) error {
	docRef := r.client.Collection(internal.USERS_COLLECTION).Doc(id)

	var update []firestore.Update

	for k, v := range u {
		update = append(update, firestore.Update{Path: k, Value: v})
	}

	_, err := docRef.Update(ctx, update)
	if err != nil {
		return err
	}
	return nil
}
