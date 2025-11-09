package services

import (
	"memora/internal/firebase"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
)

type ServiceDeps struct {
	UserRepo firebase.UserRepository
	CardRepo firebase.CardRepository
	DeckRepo firebase.DeckRepository
	AuthRepo firebase.FirebaseAuth
	Redis    *redis.Client
	Validate *validator.Validate
}

// Services groups all service instances.
type Services struct {
	Users *UserService
	Decks *DeckService
	Auth  *AuthService
}

// NewServices creates a new Services struct with the provided repositories and validator.
func NewServices(repos *firebase.Repositories, validate *validator.Validate) *Services {
	deps := &ServiceDeps{
		UserRepo: repos.User,
		CardRepo: repos.Card,
		DeckRepo: repos.Deck,
		AuthRepo: repos.Auth,
		Redis:    repos.Redis,
		Validate: validate,
	}

	return &Services{
		Users: NewUserService(deps),
		Decks: NewDeckService(deps),
		Auth:  NewAuthService(deps),
	}
}
