package services

import "memora/internal/firebase"

type DeckService struct {
	repo firebase.DeckRepository
}

func NewDeckService(repo firebase.DeckRepository) *DeckService {
	return &DeckService{repo: repo}
}