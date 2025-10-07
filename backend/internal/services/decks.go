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

type DeckService struct {
	repo     firebase.DeckRepository
	validate *validator.Validate
}

func NewDeckService(repo firebase.DeckRepository, validate *validator.Validate) *DeckService {
	return &DeckService{repo: repo, validate: validate}
}

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

func (s *DeckService) GetOneDeck(ctx context.Context, id string) (models.DeckResponse, error) {
	var response models.DeckResponse
	deck, err := s.repo.GetOneDeck(ctx, id)
	if err != nil {
		return response, err
	}

	var cards []models.Card
	for _, ref := range deck.Cards {
		snap, err := ref.Get(ctx)
		if err != nil {
			return models.DeckResponse{}, err
		}

		card, err := getCardStructFromData(snap.Data(), errors.ErrInvalidCard)
		if err != nil {
			return models.DeckResponse{}, err
		}
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

	if err := s.repo.UpdateDeck(ctx, updateMap, deckID); err != nil {
		return models.DeckResponse{}, err
	}

	return s.GetOneDeck(ctx, deckID)
}

func (s *DeckService) DeleteDeck(
	ctx context.Context,
	id string,
) error {
	return s.repo.DeleteDeck(ctx, id)
}

func getCardStructFromData(m map[string]any, errorOnFail error) (models.Card, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return GetCardStruct(data, errorOnFail)
}

func (s *DeckService) UpdateEmailsInDeck(
	ctx context.Context,
	deckID string,
	emails models.UpdateDeckEmails,
) (models.DeckResponse, error) {
	var err error

	if err := s.validate.Struct(emails); err != nil {
		return models.DeckResponse{}, errors.ErrInvalidDeck
	}

	switch emails.Opp {
	case utils.OPP_ADD:
		err = s.repo.AddEmailsToShared(ctx, deckID, emails.Emails)
	case utils.OPP_REMOVE:
		err = s.repo.RemoveEmailsFromShared(ctx, deckID, emails.Emails)
	}
	if err != nil {
		log.Println(err)
		return models.DeckResponse{}, err
	}
	return s.GetOneDeck(ctx, deckID)
}

func (s *DeckService) UpdateCardsInDeck(
	ctx context.Context,
	deckID string,
	cards models.UpdateDeckCards,
) (models.DeckResponse, error) {
	var err error

	if err := s.validate.Struct(cards); err != nil {
		return models.DeckResponse{}, errors.ErrInvalidDeck
	}

	switch cards.Opp {
	case utils.OPP_ADD:
		err = s.repo.AddCardsToDeck(ctx, deckID, cards.Cards)
	case utils.OPP_REMOVE:
		err = s.repo.RemoveCardsFromDeck(ctx, deckID, cards.Cards)
	}
	if err != nil {
		return models.DeckResponse{}, err
	}

	return s.GetOneDeck(ctx, deckID)
}
