package services

import (
	"context"
	"encoding/json"
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
	utils.MULTIPLE_CHOICE_CARD: func() models.Card { return &models.MutlipleChoiceCard{} },
	utils.ORDERED_CARD:         func() models.Card { return &models.OrderedCard{} },
}

type CardService struct {
	repo     firebase.CardRepository
	validate *validator.Validate
}

func NewCardService(repo firebase.CardRepository, validate *validator.Validate) *CardService {
	return &CardService{repo: repo, validate: validate}
}

func (s *CardService) CreateCard(ctx context.Context, rawData []byte) (string, error) {
	card, err := getCardStruct(rawData, errors.ErrInvalidCard)
	if err != nil {
		return "", err
	}

	if err := s.validate.Struct(card); err != nil {
		return "", errors.ErrInvalidCard
	}

	return s.repo.CreateCard(ctx, card)
}

func getCardStruct(data []byte, errorOnFail error) (models.Card, error) {
	var cardType models.CardType

	if err := json.Unmarshal(data, &cardType); err != nil {
		return nil, errors.ErrInvalidUser
	}

	factory, ok := cardRegistry[cardType.Type]
	if !ok {
		return nil, errorOnFail
	}

	card := factory()

	if err := json.Unmarshal(data, card); err != nil {
		return nil, err
	}

	return card, nil
}
