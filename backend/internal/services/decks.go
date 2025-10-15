package services

import (
	"context"
	"memora/internal/errors"
	"memora/internal/firebase"
	"memora/internal/models"
	"memora/internal/utils"
	"slices"

	"github.com/go-playground/validator/v10"
)

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

func (s *DeckService) CheckIfUserCanAccessDeck(
	ctx context.Context,
	deckID, userID, userEmail string,
) (bool, error) {
	deck, err := s.repo.GetOneDeck(ctx, deckID)
	if err != nil {
		return false, err
	}

	if deck.OwnerID == userID || slices.Contains(deck.SharedEmails, userEmail) {
		return true, nil
	}

	return false, nil
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
func (s *DeckService) GetOneDeck(
	ctx context.Context,
	id string,
) (models.DeckResponse, error) {
	// Fetch the deck data from the repository
	deck, err := s.repo.GetOneDeck(ctx, id)
	if err != nil {
		return models.DeckResponse{}, err
	}

	cards, err := s.Cards.GetCardsInDeck(ctx, id)
	if err != nil {
		return models.DeckResponse{}, err
	}

	return models.DeckResponse{
		ID:           id,
		Title:        deck.Title,
		OwnerID:      deck.OwnerID,
		SharedEmails: deck.SharedEmails,
		Cards:        cards,
	}, nil
}

func (s *DeckService) GetCardInDeck(
	ctx context.Context,
	deckID, cardID string,
) (models.Card, error) {
	return s.Cards.GetCardInDeck(ctx, deckID, cardID)
}

func (s *DeckService) AddCardToDeck(
	ctx context.Context,
	deckID string,
	rawData []byte,
) (models.DeckResponse, error) {
	if err := s.Cards.CreateCard(
		ctx,
		rawData,
		deckID,
	); err != nil {
		return models.DeckResponse{}, err
	}

	return s.GetOneDeck(ctx, deckID)
}

// UpdateDeck updates an existing deck identified by its ID with the provided data.
// Validates the updated deck and returns the updated deck or an error if the operation fails.
func (s *DeckService) UpdateDeck(
	ctx context.Context,
	deckID string,
	update models.UpdateDeck,
) (models.DeckResponse, error) {
	if err := s.validate.Struct(update); err != nil {
		return models.DeckResponse{}, errors.ErrInvalidDeck
	}

	updateMap, err := utils.StructToUpdate(update)
	if err != nil {
		return models.DeckResponse{}, err
	}

	// Perform the update in the repository
	if err := s.repo.UpdateDeck(ctx, updateMap, deckID); err != nil {
		return models.DeckResponse{}, err
	}

	return s.GetOneDeck(ctx, deckID)
}

func (s *DeckService) UpdateCardInDeck(
	ctx context.Context,
	deckID, cardID string,
	rawData []byte,
) (models.DeckResponse, error) {
	err := s.Cards.UpdateCard(ctx, rawData, deckID, cardID)
	if err != nil {
		return models.DeckResponse{}, err
	}

	return s.GetOneDeck(ctx, deckID)
}

// UpdateEmailsInDeck updates the shared emails of a deck based on the provided operation (add or remove).
// Validates the input and returns the updated deck or an error if the operation fails.
func (s *DeckService) UpdateEmailsInDeck(
	ctx context.Context,
	deckID string,
	emails models.UpdateDeckEmails,
) (models.DeckResponse, error) {
	var err error

	// Validate the input struct
	if err := s.validate.Struct(emails); err != nil {
		return models.DeckResponse{}, errors.ErrInvalidDeck
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
		return models.DeckResponse{}, err
	}

	// Fetch and return the updated deck
	return s.GetOneDeck(ctx, deckID)
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
