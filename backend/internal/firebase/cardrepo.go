package firebase

import (
	"context"
	"log"
	"memora/internal/config"
	"memora/internal/errors"
	"memora/internal/utils"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// CardStore defines methods for storing, retrieving and updating cards.
type CardRepository interface {
	// CreateCard adds a new card into firestore.
	// Error on fail, returns the ID if succesfull
	CreateCard(ctx context.Context, card any, deckID string) error

	GetCardsInDeck(ctx context.Context, deckID string) ([]map[string]interface{}, error)

	// GetCard returns the raw data for a card for a given ID.
	// Error on fail, or if ID is not valid
	GetCard(ctx context.Context, id string) (map[string]any, error)

	// UpdateCard updates an existing card in firestore.
	// Error on fail or if the ID is not valid, nil on success
	UpdateCard(ctx context.Context, firestoreUpdates []firestore.Update, id string) error

	// DeleteCard deletes an existing card in firestore.
	// Error on fail or if the ID is not valid, nil on success
	DeleteCard(ctx context.Context, id string) error
}

// FirestoreCardRepo holds the connection to the database
type FirestoreCardRepo struct {
	client *firestore.Client
}

// NewFirestoreCardRepo creates and returns a pointer to the repository
func NewFirestoreCardRepo(client *firestore.Client) *FirestoreCardRepo {
	return &FirestoreCardRepo{client: client}
}

func (r *FirestoreCardRepo) GetCardsInDeck(ctx context.Context, deckID string) ([]map[string]interface{}, error) {
	cardsColl := r.client.Collection(config.DecksCollection).
		Doc(deckID).
		Collection(config.CardsCollection)
	iter := cardsColl.Documents(ctx)
	defer iter.Stop()

	var result []map[string]interface{}

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		data := doc.Data()
		data["id"] = doc.Ref.ID

		result = append(result, data)
	}

	return result, nil
}

// GetCard takes a context and id, and returns the raw data for a card
// error if the card can not be fetched
func (r *FirestoreCardRepo) GetCard(ctx context.Context, id string) (map[string]any, error) {
	doc, err := r.client.Collection(config.CardsCollection).Doc(id).Get(ctx)
	if err != nil {
		return nil, errors.ErrInvalidId
	}

	return doc.Data(), nil
}

// CreateCard takes a context and a card, adds it to the database, and
// returns the created card or an error if the operation fails.
func (r *FirestoreCardRepo) CreateCard(ctx context.Context, card any, deckID string) error {
	log.Printf("Repo: adding card to deck %s: %+v", deckID, card)
	_, _, err := r.client.Collection(config.DecksCollection).
		Doc(deckID).
		Collection(config.CardsCollection).
		Add(ctx, card)

	log.Println(err)
	return err
}

// UpdateCard takes a context, an update payload, and an ID, and updates
// the corresponding card in the database. It returns an error if the update
// fails or the card cannot be found
func (r *FirestoreCardRepo) UpdateCard(
	ctx context.Context,
	firestoreUpdates []firestore.Update,
	id string,
) error {
	return utils.UpdateDocumentInDB(
		r.client,
		ctx,
		config.CardsCollection,
		id,
		firestoreUpdates,
	)
}

// DeleteCard takes a context and ID, and deletes teh corresponding
// card in teh database. It returns an error if the delete fails,
// or if it cannot be found
func (r *FirestoreCardRepo) DeleteCard(ctx context.Context, id string) error {
	return utils.DeleteDocumentInDB(
		r.client,
		ctx,
		config.CardsCollection,
		id,
	)
}
