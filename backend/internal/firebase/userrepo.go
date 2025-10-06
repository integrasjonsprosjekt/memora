package firebase

import (
	"context"
	"memora/internal/config"
	"memora/internal/errors"
	"memora/internal/models"
	"memora/internal/utils"

	"cloud.google.com/go/firestore"
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

func (r *FirestoreUserRepo) AddUser(ctx context.Context, user models.CreateUser) (string, error) {
	exists, err := utils.UserExistsByEmail(r.client, ctx, user.Email)
	if err != nil {
		return "", err
	}
	if exists {
		return "", errors.ErrInvalidEmailPresent
	}

	return utils.AddToDB(r.client, ctx, config.UsersCollection, user)
}

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

func (r *FirestoreUserRepo) DeleteUser(ctx context.Context, id string) error {
	return utils.DeleteDocumentInDB(
		r.client,
		ctx,
		config.UsersCollection,
		id,
	)
}
