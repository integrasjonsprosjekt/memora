package services

import (
	"memora/internal/firebase"

	"github.com/go-playground/validator/v10"
)

type CardService struct {
	repo     firebase.CardRepository
	validate *validator.Validate
}

func NewCardService(repo firebase.CardRepository, validate *validator.Validate) *CardService {
	return &CardService{repo: repo, validate: validate}
}
