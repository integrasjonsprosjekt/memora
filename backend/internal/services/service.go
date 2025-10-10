package services

import (
	"memora/internal/firebase"

	"github.com/go-playground/validator/v10"
)

// Services groups all service instances.
type Services struct {
	Users *UserService
	Decks *DeckService
}

// NewServices creates a new Services struct with the provided repositories and validator.
func NewServices(
	repos *firebase.Repositories,
	validate *validator.Validate,
) *Services {

	return &Services{
		Users: NewUserService(
			repos.User,
			validate,
		),
		Decks: NewDeckService(
			repos.Deck,
			validate,
			NewCardService(repos.Card, validate),
		),
	}
}
