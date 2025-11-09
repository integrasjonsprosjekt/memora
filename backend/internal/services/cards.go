package services

import (
	"context"
	"encoding/json"
	"fmt"
	"memora/internal/errors"
	"memora/internal/firebase"
	"memora/internal/models"
	"memora/internal/utils"
	"strconv"
	"time"

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
func NewCardService(
	repo firebase.CardRepository,
	validate *validator.Validate,
) *CardService {
	return &CardService{
		repo:     repo,
		validate: validate,
	}
}

// GetCard retrieves a card by its ID.
// Returns the card or an error if the operation fails.
func (s *CardService) GetCardInDeck(
	ctx context.Context,
	deckID, cardID string,
) (models.Card, error) {
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

// GetCardsInDeck retrieves all cards in a specified deck with cursor-based pagination.
// cursor is the ID of the last card from the previous page (empty string for first page)
// Returns a list of cards or an error if the operation fails.
func (s *CardService) GetCardsInDeck(
	ctx context.Context,
	deckID, limit_str string,
	cursor string,
) ([]models.Card, bool, error) {
	limit := utils.ParseLimit(limit_str)
	// Fetch raw card documents from the repository
	docs, hasMore, err := s.repo.GetCardsInDeck(ctx, deckID, limit, cursor)
	if err != nil {
		return nil, false, err
	}

	// Convert each document to the appropriate card struct
	// based on its type
	var cards []models.Card
	for _, doc := range docs {
		raw, err := json.Marshal(doc)
		if err != nil {
			return nil, false, err
		}

		card, err := GetCardStruct(raw, fmt.Errorf("internal server error"))
		if err != nil {
			return nil, false, err
		}

		card.SetID(doc["id"].(string))
		cards = append(cards, card)
	}

	return cards, hasMore, nil
}

// CreateCard creates a new card from the provided raw JSON data.
// Validates the card and returns its ID or an error if the operation fails.
func (s *CardService) CreateCard(
	ctx context.Context,
	rawData []byte,
	deckID string,
) (string, error) {
	// Parse the raw data to determine the card type and unmarshal into the correct struct
	card, err := GetCardStruct(rawData, errors.ErrInvalidCard)
	if err != nil {
		return "", err
	}

	if err := s.validate.Struct(card); err != nil {
		return "", errors.ErrInvalidCard
	}

	return s.repo.CreateCard(ctx, card, deckID)
}

// UpdateCard updates an existing card identified by its ID with the provided raw JSON data.
// Validates the updated card and returns the updated card or an error if the operation fails.
func (s CardService) UpdateCard(
	ctx context.Context,
	rawData []byte,
	deckID, cardID string,
) error {
	card, err := GetCardStruct(rawData, errors.ErrInvalidCard)
	if err != nil {
		return err
	}

	originalCard, err := s.repo.GetCardInDeck(ctx, deckID, cardID)
	if err != nil {
		return err
	}

	// Ensure the card type is not being changed
	t, ok := originalCard["type"].(string)
	if !ok {
		return fmt.Errorf("internal server error")
	}

	if t != card.GetType() {
		return errors.ErrInvalidCard
	}

	if err := s.validate.Struct(card); err != nil {
		return errors.ErrInvalidCard
	}

	// Convert the updated card struct to firestore updates
	update, err := utils.StructToUpdate(card)
	if err != nil {
		return errors.ErrInvalidCard
	}

	// Perform the update in the repository
	return s.repo.UpdateCard(ctx, update, deckID, cardID)
}

// DeleteCard deletes a card by its ID.
// Returns an error if the operation fails or the card is not found.
func (s *CardService) DeleteCard(
	ctx context.Context,
	deckID, cardID string,
) error {
	return s.repo.DeleteCard(ctx, deckID, cardID)
}

// GetCardStruct takes a byte array and an error to return if the type is not found.
// It returns a card struct of the appropriate type based on the "type" field in the JSON data.
func GetCardStruct(
	data []byte,
	errorOnFail error,
) (models.Card, error) {
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

func (s *CardService) GetCardProgress(
	ctx context.Context,
	deckID, cardID, userID string,
) (models.CardProgress, error) {
	return s.repo.GetCardProgress(ctx, deckID, cardID, userID)
}

func (s *CardService) UpdateCardProgress(
	ctx context.Context,
	deckID, cardID, userID string,
	rating models.CardRating,
) error {
	if err := s.validate.Struct(rating); err != nil {
		return errors.ErrInvalidUser
	}

	progress, err := s.GetCardProgress(ctx, deckID, cardID, userID)
	if err != nil {
		if err == errors.ErrInvalidId {
			progress = models.CardProgress{
				EaseFactor:   2500,
				Reps:         1,
				Lapses:       0,
				Interval:     0,
				LastReviewed: time.Time{},
				Due:          time.Time{},
			}
		} else {
			return err
		}
	}

	now := time.Now()
	easeFactor := progress.EaseFactor
	reps := progress.Reps
	lapses := progress.Lapses
	interval := float64(progress.Interval)

	switch rating.Rating {
	case "again":
		easeFactor -= 200
		lapses += 1
		reps += 1
		interval = 1.0
	case "hard":
		easeFactor -= 150
		reps += 1
		interval *= 1.2
	case "good":
		reps += 1
		interval *= 1.5
	case "easy":
		easeFactor += 150
		reps += 1
		interval *= 2.0
	}

	if easeFactor < 1300 {
		easeFactor = 1300
	}

	if easeFactor > 3000 {
		easeFactor = 3000
	}

	progress.EaseFactor = easeFactor
	progress.Reps = reps
	progress.Lapses = lapses
	progress.Interval = interval
	progress.LastReviewed = now
	progress.Due = now.Add(time.Duration(interval*24) * time.Hour)

	return s.repo.UpdateProgress(ctx, deckID, cardID, userID, progress)
}

func (s *CardService) GetDueCardsInDeck(
	ctx context.Context,
	deckID, userID string,
	limit, cursor string,
) ([]models.Card, string, bool, error) {
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return nil, "", false, err
	}

	docs, nextCursor, hasMore, err := s.repo.GetDueCardsInDeck(
		ctx,
		deckID,
		userID,
		limitInt,
		cursor,
	)
	if err != nil {
		return nil, "", false, err
	}

	var cards []models.Card
	for _, doc := range docs {
		raw, err := json.Marshal(doc)
		if err != nil {
			return nil, "", false, err
		}

		card, err := GetCardStruct(raw, fmt.Errorf("internal server error"))
		if err != nil {
			return nil, "", false, err
		}

		card.SetID(doc["id"].(string))
		cards = append(cards, card)
	}

	return cards, nextCursor, hasMore, nil
}
