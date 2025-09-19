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
	id, err := s.repo.AddDeck(ctx, deck)
	if err != nil {
		return "", err
	}

	return id, nil
}
