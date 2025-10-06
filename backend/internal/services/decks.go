package services

import (
	"context"
	"encoding/json"
	"memora/internal/errors"
	"memora/internal/firebase"
	"memora/internal/models"
	"memora/internal/utils"

	"cloud.google.com/go/firestore"
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
			return models.DeckResponse{}, nil
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
	var err error

	if err := s.validate.Struct(update); err != nil {
		return models.DeckResponse{}, errors.ErrInvalidDeck
	}
	if len(update.Cards) > 0 && update.Operation != "" {
		for _, id := range update.Cards {
			switch update.Operation {
			case "add":
				err = s.repo.AddCardToDeck(ctx, deckID, id)
			case "remove":
				err = s.repo.RemoveCardFromDeck(ctx, deckID, id)
			default:
				return models.DeckResponse{}, errors.ErrInvalidDeck
			}

			if err != nil {
				return models.DeckResponse{}, err
			}
		}
	}

	updates, err := utils.StructToUpdate(update)
	if err != nil {
		return models.DeckResponse{}, err
	}

	filterd := make([]firestore.Update, 0, len(updates))
	for _, u := range updates {
		switch u.Path {
		case "cards", "operation":
			continue
		default:
			filterd = append(filterd, u)
		}
	}

	if len(filterd) > 0 {
		if err := s.repo.UpdateDeck(ctx, filterd, deckID); err != nil {
			return models.DeckResponse{}, err
		}
	}

	return s.GetOneDeck(ctx, deckID)
}

func getCardStructFromData(m map[string]any, errorOnFail error) (models.Card, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return GetCardStruct(data, errorOnFail)
}
