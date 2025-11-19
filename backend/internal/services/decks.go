package services

import (
	"context"
	"memora/internal/errors"
	"memora/internal/firebase"
	"memora/internal/models"
	"memora/internal/utils"
	"slices"
	"sync"

	"github.com/go-playground/validator/v10"
)

// Default filter for all fields, used when updating a deck
const defaultFilterDecks = "title,owner_id,shared_emails"

// DeckService provides methods for managing decks.
type DeckService struct {
	repo     firebase.DeckRepository
	validate *validator.Validate
	cache    *CacheService
	Cards    *CardService
}

// NewDeckService creates a new instance of DeckService.
func NewDeckService(
	deps *ServiceDeps,
) *DeckService {
	return &DeckService{
		repo:     deps.DeckRepo,
		validate: deps.Validate,
		cache:    deps.Cache,
		Cards:    NewCardService(deps),
	}
}

func (s *DeckService) GetCardsInDeck(
	ctx context.Context,
	deckID, limit_str, cursor string,
) ([]models.Card, bool, error) {
	return s.Cards.GetCardsInDeck(ctx, deckID, limit_str, cursor)
}

func (s *DeckService) CheckIfUserCanAccessDeck(
	ctx context.Context,
	deckID, userID, userEmail string,
) (bool, error) {
	deck, err := s.repo.GetOneDeck(ctx, deckID, []string{"owner_id", "shared_emails"})
	if err != nil {
		return false, err
	}

	if deck.OwnerID == userID || slices.Contains(deck.SharedEmails, userEmail) {
		return true, nil
	}

	return false, nil
}

func (s *DeckService) UserOwnsDeck(
	ctx context.Context,
	deckID, userID string,
) (bool, error) {
	deck, err := s.repo.GetOneDeck(ctx, deckID, []string{"owner_id"})
	if err != nil {
		return false, err
	}

	if deck.OwnerID == userID {
		return true, nil
	}

	return false, nil
}

// RegisterNewDeck creates a new deck from the provided data.
// Validates the deck and returns its ID or an error if the operation fails.
func (s *DeckService) RegisterNewDeck(
	ctx context.Context,
	deck models.CreateDeck,
	ownerEmail string,
) (string, error) {
	if err := s.validate.Struct(deck); err != nil {
		return "", errors.ErrInvalidDeck
	}

	id, err := s.repo.AddDeck(ctx, deck)
	if err != nil {
		return "", err
	}

	s.invalidateDeckCaches(id, ownerEmail, deck.SharedEmails)

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

	cacheKey := utils.DeckKey(id)
	var deck models.Deck
	err = s.cache.Get(ctx, cacheKey, &deck)
	if err == nil {
		return deck, nil
	}

	// Fetch the deck data from the repository
	deck, err = s.repo.GetOneDeck(ctx, id, filterParsed)
	if err != nil {
		return models.Deck{}, err
	}

	// Store the deck in the cache for future requests
	s.cache.SetAsync(cacheKey, deck, DeckTTL)

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
	deckID, ownerEmail string,
	update models.UpdateDeck,
) (models.Deck, error) {
	if err := s.validate.Struct(update); err != nil {
		return models.Deck{}, errors.ErrInvalidDeck
	}

	updateMap, err := utils.StructToUpdate(update)
	if err != nil {
		return models.Deck{}, err
	}

	deck, err := s.repo.GetOneDeck(ctx, deckID, []string{"shared_emails"})
	if err != nil {
		return models.Deck{}, err
	}

	// Perform the update in the repository
	if err := s.repo.UpdateDeck(ctx, updateMap, deckID); err != nil {
		return models.Deck{}, err
	}

	s.invalidateDeckCaches(deckID, ownerEmail, deck.SharedEmails)

	// Fetch and return the updated deck
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
	deckID, ownerEmail string,
	emails models.UpdateDeckEmails,
) (models.Deck, error) {
	var err error

	// Validate the input struct
	if err := s.validate.Struct(emails); err != nil {
		return models.Deck{}, errors.ErrInvalidDeck
	}

	oldDeck, err := s.repo.GetOneDeck(ctx, deckID, []string{"shared_emails"})
	if err != nil {
		return models.Deck{}, err
	}
	oldShared := oldDeck.SharedEmails

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

	newDeck, err := s.repo.GetOneDeck(ctx, deckID, []string{"shared_emails"})
	if err != nil {
		return models.Deck{}, err
	}
	newShared := newDeck.SharedEmails

	// Invalidate caches for both old and new shared emails
	allShared := append(oldShared, newShared...)
	s.invalidateDeckCaches(deckID, ownerEmail, allShared)

	// Fetch and return the updated deck
	return s.GetOneDeck(ctx, deckID, defaultFilterDecks)
}

func (s *DeckService) DeleteDeck(
	ctx context.Context,
	id, ownerEmail string,
) error {
	err := s.repo.DeleteDeck(ctx, id)
	if err != nil {
		return err
	}

	s.invalidateDeckCaches(id, ownerEmail, nil)

	return nil
}

// DeleteCardInDeck deletes a card from a deck by their IDs.
// Returns an error if the operation fails or the card is not found.
func (s *DeckService) DeleteCardInDeck(
	ctx context.Context,
	deckID, cardID string,
) error {
	return s.Cards.DeleteCard(ctx, deckID, cardID)
}

func (s *DeckService) GetCardProgress(
	ctx context.Context,
	deckID, cardID, userID string,
) (models.CardProgress, error) {
	return s.Cards.GetCardProgress(ctx, deckID, cardID, userID)
}

func (s *DeckService) UpdateCardProgress(
	ctx context.Context,
	deckID, cardID, userID string,
	rating models.CardRating,
) error {
	return s.Cards.UpdateCardProgress(ctx, deckID, cardID, userID, rating)
}

func (s *DeckService) GetDueCardsInDeck(
	ctx context.Context,
	deckID, userID string,
	limit, cursor string,
) ([]models.Card, string, bool, error) {
	return s.Cards.GetDueCardsInDeck(ctx, deckID, userID, limit, cursor)
}

func (s *DeckService) invalidateDeckCaches(deckID, ownerEmail string, sharedEmails []string) {
	ctx, cancel := context.WithTimeout(context.Background(), CacheOpTimeout)
	defer cancel()

	var wg sync.WaitGroup

	// Delete deck cache
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.cache.Delete(ctx, utils.DeckKey(deckID))
	}()

	// Invalidate owner's deck list
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.cache.Delete(ctx, utils.UserEmailDecksKey(ownerEmail))
	}()

	// Invalidate shared users' deck lists
	for _, email := range sharedEmails {
		wg.Add(1)
		go func(e string) {
			defer wg.Done()
			s.cache.Delete(ctx, utils.UserEmailDecksKey(e))
		}(email)
	}

	wg.Wait()
}
