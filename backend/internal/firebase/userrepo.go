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
	GetDecksOwned(ctx context.Context, id string) ([]models.DisplayDeck, error)
	GetDecksShared(ctx context.Context, id string) ([]models.DisplayDeck, error)
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

func (r *FirestoreUserRepo) GetDecksOwned(
	ctx context.Context,
	id string,
) ([]models.DisplayDeck, error) {
	var decks []models.DisplayDeck

	_, err := utils.GetDocumentIfExists(r.client, ctx, config.UsersCollection, id)
	if err != nil {
		return nil, err
	}

	iter := r.client.Collection(config.DecksCollection).
		Where("owner_id", "==", id).
		Documents(ctx)

	for {
		doc, err := iter.Next()

		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var deck models.DisplayDeck
		if err := doc.DataTo(&deck); err != nil {
			return nil, err
		}

		deck.ID = doc.Ref.ID
		decks = append(decks, deck)
	}

	return decks, nil
}

func (r *FirestoreUserRepo) GetDecksShared(
	ctx context.Context,
	id string,
) ([]models.DisplayDeck, error) {
	var decks []models.DisplayDeck

	user, err := utils.FetchByID[models.User](r.client, ctx, config.UsersCollection, id)
	if err != nil {
		return nil, err
	}

	iter := r.client.Collection(config.DecksCollection).
		Where("shared_emails", "array-contains", user.Email).
		Documents(ctx)

	for {
		doc, err := iter.Next()

		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var deck models.DisplayDeck
		if err := doc.DataTo(&deck); err != nil {
			return nil, err
		}

		deck.ID = doc.Ref.ID
		decks = append(decks, deck)
	}

	return decks, nil
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
