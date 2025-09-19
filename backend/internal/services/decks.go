package services

import (
	"context"
	"memora/internal/firebase"
	"memora/internal/models"
)

type DeckService struct {
	repo firebase.DeckRepository
}

func NewDeckService(repo firebase.DeckRepository) *DeckService {
	return &DeckService{repo: repo}
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
	resp, err := s.repo.GetOneDeck(ctx, id)
	if err != nil {
		return resp, err
	}

	if resp.SharedEmails == nil {
		resp.SharedEmails = []string{}
	}
	return resp, nil
}

func (s *DeckService) UpdateDeck(ctx context.Context, id string, update models.UpdateDeck) error {
	return s.repo.UpdateDeck(ctx, id, update)
}
