package services

import (
	"context"
	"memora/internal/errors"
	"memora/internal/firebase"
	"memora/internal/models"
	"memora/internal/utils"

	"github.com/go-playground/validator/v10"
)

// Default filter for all fields, used when updating a deck
const defaultFilterDecks = "title,owner_id,shared_emails"

// DeckService provides methods for managing decks.
type DeckService struct {
	repo     firebase.DeckRepository
	validate *validator.Validate
	Cards    *CardService
}

// NewDeckService creates a new instance of DeckService.
func NewDeckService(
	repo firebase.DeckRepository,
	validate *validator.Validate,
	cards *CardService,
) *DeckService {
	return &DeckService{
		repo:     repo,
		validate: validate,
		Cards:    cards,
	}
}

func (s *DeckService) GetCardsInDeck(
	ctx context.Context,
	deckID, limit_str, cursor string,
) ([]models.Card, bool, error) {
	limit, err := utils.ParseLimit(limit_str)
	if err != nil {
		return nil, false, err
	}

	return s.Cards.GetCardsInDeck(ctx, deckID, limit, cursor)
}

// RegisterNewDeck creates a new deck from the provided data.
// Validates the deck and returns its ID or an error if the operation fails.
func (s *DeckService) RegisterNewDeck(
	ctx context.Context,
	deck models.CreateDeck,
) (string, error) {
	if err := s.validate.Struct(deck); err != nil {
		return "", errors.ErrInvalidDeck
	}

	id, err := s.repo.AddDeck(ctx, deck)
	if err != nil {
		return "", err
	}

	return id, nil
}

// GetOneDeck retrieves a deck by its ID, including its cards.
// Returns the deck or an error if the operation fails.
// Filter specifies which fields to return.
func (s *DeckService) GetOneDeck(
	ctx context.Context,
	id string,
	filter string,
) (models.Deck, error) {
	filterParsed, err := utils.ParseFilter(filter)
	if err != nil {
		return models.Deck{}, err
	}

	// Fetch the deck data from the repository
	deck, err := s.repo.GetOneDeck(ctx, id, filterParsed)
	if err != nil {
		return models.Deck{}, err
	}

	return deck, nil
}

// GetCardInDeck retrieves a specific card from a deck by their IDs.
// Returns the card or an error if the operation fails.
func (s *DeckService) GetCardInDeck(
	ctx context.Context,
	deckID, cardID string,
) (models.Card, error) {
	return s.Cards.GetCardInDeck(ctx, deckID, cardID)
}

// AddCardToDeck creates a new card in the specified deck from the provided raw JSON data.
// Validates the card and returns the updated deck or an error if the operation fails.
func (s *DeckService) AddCardToDeck(
	ctx context.Context,
	deckID string,
	rawData []byte,
) (models.Card, error) {
	id, err := s.Cards.CreateCard(
		ctx,
		rawData,
		deckID,
	)
	if err != nil {
		return nil, err
	}

	return s.Cards.GetCardInDeck(ctx, deckID, id)
}

// UpdateDeck updates an existing deck identified by its ID with the provided data.
// Validates the updated deck and returns the updated deck or an error if the operation fails.
func (s *DeckService) UpdateDeck(
	ctx context.Context,
	deckID string,
	update models.UpdateDeck,
) (models.Deck, error) {
	if err := s.validate.Struct(update); err != nil {
		return models.Deck{}, errors.ErrInvalidDeck
	}

	updateMap, err := utils.StructToUpdate(update)
	if err != nil {
		return models.Deck{}, err
	}

	// Perform the update in the repository
	if err := s.repo.UpdateDeck(ctx, updateMap, deckID); err != nil {
		return models.Deck{}, err
	}

	return s.GetOneDeck(ctx, deckID, defaultFilterDecks)
}

// UpdateCardInDeck updates an existing card in a deck with the provided raw JSON data.
// Validates the updated card and returns the updated deck or an error if the operation fails.
func (s *DeckService) UpdateCardInDeck(
	ctx context.Context,
	deckID, cardID string,
	rawData []byte,
) (models.Card, error) {
	err := s.Cards.UpdateCard(ctx, rawData, deckID, cardID)
	if err != nil {
		return nil, err
	}

	return s.GetCardInDeck(ctx, deckID, cardID)
}

// UpdateEmailsInDeck updates the shared emails of a deck based on the provided operation (add or remove).
// Validates the input and returns the updated deck or an error if the operation fails.
func (s *DeckService) UpdateEmailsInDeck(
	ctx context.Context,
	deckID string,
	emails models.UpdateDeckEmails,
) (models.Deck, error) {
	var err error

	// Validate the input struct
	if err := s.validate.Struct(emails); err != nil {
		return models.Deck{}, errors.ErrInvalidDeck
	}

	// Perform the appropriate operation based on the Opp field
	switch emails.Opp {
	case utils.OPP_ADD:
		// Add emails to the deck's shared emails
		err = s.repo.AddEmailsToShared(ctx, deckID, emails.Emails)
	case utils.OPP_REMOVE:
		// Remove emails from the deck's shared emails
		err = s.repo.RemoveEmailsFromShared(ctx, deckID, emails.Emails)
	}
	if err != nil {
		return models.Deck{}, err
	}

	// Fetch and return the updated deck
	return s.GetOneDeck(ctx, deckID, defaultFilterDecks)
}

// DeleteDeck deletes a deck by its ID.
// Returns an error if the operation fails or the deck is not found.
func (s *DeckService) DeleteDeck(
	ctx context.Context,
	id string,
) error {
	return s.repo.DeleteDeck(ctx, id)
}

// DeleteCardInDeck deletes a card from a deck by their IDs.
// Returns an error if the operation fails or the card is not found.
func (s *DeckService) DeleteCardInDeck(
	ctx context.Context,
	deckID, cardID string,
) error {
	return s.Cards.DeleteCard(ctx, deckID, cardID)
}
