package services

import (
	"context"
	"encoding/json"
	"fmt"
	"memora/internal/errors"
	"memora/internal/firebase"
	"memora/internal/models"
	"memora/internal/utils"

	"github.com/go-playground/validator/v10"
)

// Used to get the type of card based on request
var cardRegistry = map[string]func() models.Card{
	utils.BLANKS_CARD:          func() models.Card { return &models.BlanksCard{} },
	utils.FRONT_BACK_CARD:      func() models.Card { return &models.FrontBackCard{} },
	utils.MULTIPLE_CHOICE_CARD: func() models.Card { return &models.MultipleChoiceCard{} },
	utils.ORDERED_CARD:         func() models.Card { return &models.OrderedCard{} },
}

// CardService provides methods for managing cards.
type CardService struct {
	repo     firebase.CardRepository
	validate *validator.Validate
}

// NewCardService creates a new instance of CardService.
func NewCardService(repo firebase.CardRepository, validate *validator.Validate) *CardService {
	return &CardService{repo: repo, validate: validate}
}

// GetCard retrieves a card by its ID.
// Returns the card or an error if the operation fails.
func (s *CardService) GetCardInDeck(ctx context.Context, deckID, cardID string) (models.Card, error) {
	doc, err := s.repo.GetCardInDeck(ctx, deckID, cardID)
	if err != nil {
		return nil, err
	}

	// Convert the document to JSON
	raw, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}

	// Convert the JSON to the appropriate card struct based on its type
	card, err := GetCardStruct(raw, fmt.Errorf("internal server error"))
	if err != nil {
		return nil, err
	}

	card.SetID(cardID)

	return card, nil
}

func (s *CardService) GetCardsInDeck(ctx context.Context, deckID string) ([]models.Card, error) {
	docs, err := s.repo.GetCardsInDeck(ctx, deckID)
	if err != nil {
		return nil, err
	}

	var cards []models.Card
	for _, doc := range docs {
		raw, err := json.Marshal(doc)
		if err != nil {
			return nil, err
		}

		card, err := GetCardStruct(raw, fmt.Errorf("internal server error"))
		if err != nil {
			return nil, err
		}

		card.SetID(doc["id"].(string))

		cards = append(cards, card)
	}

	return cards, nil
}

// CreateCard creates a new card from the provided raw JSON data.
// Validates the card and returns its ID or an error if the operation fails.
func (s *CardService) CreateCard(ctx context.Context, rawData []byte, deckID string) error {
	// Parse the raw data to determine the card type and unmarshal into the correct struct
	card, err := GetCardStruct(rawData, errors.ErrInvalidCard)
	if err != nil {
		return err
	}

	if err := s.validate.Struct(card); err != nil {
		return errors.ErrInvalidCard
	}

	return s.repo.CreateCard(ctx, card, deckID)
}

/*
// UpdateCard updates an existing card identified by its ID with the provided raw JSON data.
// Validates the updated card and returns the updated card or an error if the operation fails.
func (s CardService) UpdateCard(ctx context.Context, rawData []byte, id string) (any, error) {
	card, err := GetCardStruct(rawData, errors.ErrInvalidCard)
	if err != nil {
		return nil, err
	}

	originalCard, err := s.repo.GetCard(ctx, id)
	if err != nil {
		return nil, err
	}

	// Ensure the card type is not being changed
	t, ok := originalCard["type"].(string)
	if !ok {
		return nil, fmt.Errorf("internal server error")
	}

	if t != card.GetType() {
		return nil, errors.ErrInvalidCard
	}

	if err := s.validate.Struct(card); err != nil {
		return nil, errors.ErrInvalidCard
	}

	// Convert the updated card struct to firestore updates
	update, err := utils.StructToUpdate(card)
	if err != nil {
		return nil, errors.ErrInvalidCard
	}

	// Perform the update in the repository
	err = s.repo.UpdateCard(ctx, update, id)
	if err != nil {
		return nil, err
	}

	// Fetch and return the updated card
	returnCard, err := s.GetCard(ctx, id)
	if err != nil {
		return nil, err
	}

	return returnCard, nil
}*/

// DeleteCard deletes a card by its ID.
// Returns an error if the operation fails or the card is not found.
func (s *CardService) DeleteCard(ctx context.Context, deckID, cardID string) error {
	return s.repo.DeleteCard(ctx, deckID, cardID)
}

// GetCardStruct takes a byte array and an error to return if the type is not found.
// It returns a card struct of the appropriate type based on the "type" field in the JSON data.
func GetCardStruct(data []byte, errorOnFail error) (models.Card, error) {
	var cardType models.CardType

	// First, unmarshal to get the card type
	if err := json.Unmarshal(data, &cardType); err != nil {
		return nil, errors.ErrInvalidUser
	}

	// Secondly, lookup the card type in the registry and create a new instance
	factory, ok := cardRegistry[cardType.Type]
	if !ok {
		return nil, errorOnFail
	}

	card := factory()

	// Thirdly, unmarshal the JSON data into the specific card struct
	if err := json.Unmarshal(data, card); err != nil {
		return nil, err
	}

	return card, nil
}
