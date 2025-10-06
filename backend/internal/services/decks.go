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
	if err := s.validate.Struct(update); err != nil {
		return models.DeckResponse{}, errors.ErrInvalidDeck
	}

	if err := s.updateCardsInDeck(
		ctx,
		update.OppCards,
		deckID,
		update.Cards,
	); err != nil {
		return models.DeckResponse{}, err
	}

	if err := s.updateEmailsInDeck(
		ctx,
		update.OppEmails,
		deckID,
		update.Emails,
	); err != nil {
		return models.DeckResponse{}, err
	}

	updates, err := utils.StructToUpdate(update)
	if err != nil {
		return models.DeckResponse{}, err
	}

	filterd := s.filterOutArray(updates, "opp_cards", "cards", "opp_emails", "emails")

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

func (s *DeckService) updateEmailsInDeck(
	ctx context.Context,
	opp, deckID string,
	emails []string,
) error {
	var err error
	for _, email := range emails {
		switch opp {
		case utils.OPP_ADD:
			err = s.repo.AddEmailToShared(ctx, email, deckID)
		case utils.OPP_REMOVE:
			err = s.repo.RemoveEmailFromShared(ctx, email, deckID)
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func (s *DeckService) updateCardsInDeck(
	ctx context.Context,
	opp, deckID string,
	field []string,
) error {
	var err error
	for _, id := range field {
		switch opp {
		case utils.OPP_ADD:
			err = s.repo.AddCardToDeck(ctx, deckID, id)
		case utils.OPP_REMOVE:
			err = s.repo.RemoveCardFromDeck(ctx, deckID, id)
		default:
			return errors.ErrInvalidDeck
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *DeckService) filterOutArray(
	updates []firestore.Update,
	args ...any,
) []firestore.Update {
	filterd := make([]firestore.Update, 0, len(updates))
	skipFields := map[string]struct{}{}
	for _, arg := range args {
		if field, ok := arg.(string); ok {
			skipFields[field] = struct{}{}
		}
	}

	for _, update := range updates {
		if _, skip := skipFields[update.Path]; skip {
			continue
		}
		filterd = append(filterd, update)
	}

	return filterd
}
