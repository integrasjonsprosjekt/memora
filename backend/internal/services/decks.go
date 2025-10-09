package services

import (
	"context"
	"encoding/json"
	"log"
	"memora/internal/errors"
	"memora/internal/firebase"
	"memora/internal/models"
	"memora/internal/utils"

	"github.com/go-playground/validator/v10"
)

// DeckService provides methods for managing decks.
type DeckService struct {
	repo     firebase.DeckRepository
	validate *validator.Validate
}

// NewDeckService creates a new instance of DeckService.
func NewDeckService(repo firebase.DeckRepository, validate *validator.Validate) *DeckService {
	return &DeckService{repo: repo, validate: validate}
}

// RegisterNewDeck creates a new deck from the provided data.
// Validates the deck and returns its ID or an error if the operation fails.
func (s *DeckService) RegisterNewDeck(ctx context.Context, deck models.CreateDeck) (string, error) {
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
func (s *DeckService) GetOneDeck(ctx context.Context, id string) (models.DeckResponse, error) {
	var response models.DeckResponse

	// Fetch the deck data from the repository
	deck, err := s.repo.GetOneDeck(ctx, id)
	if err != nil {
		return response, err
	}

	// Fetch all cards associated with the deck
	var cards []models.Card
	for _, ref := range deck.Cards {
		snap, err := ref.Get(ctx)
		if err != nil {
			return models.DeckResponse{}, err
		}

		// Convert the snapshot data to a card struct
		card, err := getCardStructFromData(snap.Data(), errors.ErrInvalidCard)
		if err != nil {
			return models.DeckResponse{}, err
		}

		// Set the card ID from the document reference
		card.SetID(snap.Ref.ID)
		cards = append(cards, card)
	}

	return models.DeckResponse{
		ID:           id,
		Title:        deck.Title,
		OwnerID:      deck.OwnerID,
		SharedEmails: deck.SharedEmails,
		Cards:        cards,
	}, nil
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

// DeleteDeck deletes a deck by its ID.
// Returns an error if the operation fails or the deck is not found.
func (s *DeckService) DeleteDeck(
	ctx context.Context,
	id string,
) error {
	return s.repo.DeleteDeck(ctx, id)
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
		log.Println(err)
		return models.DeckResponse{}, err
	}

	// Fetch and return the updated deck
	return s.GetOneDeck(ctx, deckID)
}

// UpdateCardsInDeck updates the cards of a deck based on the provided operation (add or remove).
// Validates the input and returns the updated deck or an error if the operation fails.
func (s *DeckService) UpdateCardsInDeck(
	ctx context.Context,
	deckID string,
	cards models.UpdateDeckCards,
) (models.DeckResponse, error) {
	var err error

	// Validate the input struct
	if err := s.validate.Struct(cards); err != nil {
		return models.DeckResponse{}, errors.ErrInvalidDeck
	}

	// Perform the appropriate operation based on the Opp field
	switch cards.Opp {
	case utils.OPP_ADD:
		// Add cards to the deck
		err = s.repo.AddCardsToDeck(ctx, deckID, cards.Cards)
	case utils.OPP_REMOVE:
		// Remove cards from the deck
		err = s.repo.RemoveCardsFromDeck(ctx, deckID, cards.Cards)
	}
	if err != nil {
		return models.DeckResponse{}, err
	}

	// Fetch and return the updated deck
	return s.GetOneDeck(ctx, deckID)
}

// getCardStructFromData converts a map of raw data to a Card struct based on its type.
func getCardStructFromData(m map[string]any, errorOnFail error) (models.Card, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return GetCardStruct(data, errorOnFail)
}