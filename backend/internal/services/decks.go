package services

import (
	"context"
	"encoding/json"
	"memora/internal/errors"
	"memora/internal/firebase"
	"memora/internal/models"

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
	if deck.SharedEmails == nil {
		deck.SharedEmails = []string{}
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
			return models.DeckResponse{}, nil
		}
		card.SetID(snap.Ref.ID)

		cards = append(cards, card)
	}

	return models.DeckResponse{
		ID:           id,
		OwnerID:      deck.OwnerID,
		SharedEmails: deck.SharedEmails,
		Cards:        cards,
	}, nil
}

func (s *DeckService) UpdateDeck(ctx context.Context, id string, update models.UpdateDeck) error {
	return s.repo.UpdateDeck(ctx, id, update)
}

func getCardStructFromData(m map[string]any, errorOnFail error) (models.Card, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return GetCardStruct(data, errorOnFail)
}
