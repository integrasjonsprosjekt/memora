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
	Cache    *CacheService
	Validate *validator.Validate
}

// Services groups all service instances.
type Services struct {
	Users *UserService
	Decks *DeckService
	Auth  *AuthService
	Rdb   *redis.Client
}

// NewServices creates a new Services struct with the provided repositories and validator.
func NewServices(repos *firebase.Repositories, validate *validator.Validate, rdb *redis.Client) *Services {
	deps := &ServiceDeps{
		UserRepo: repos.User,
		CardRepo: repos.Card,
		DeckRepo: repos.Deck,
		AuthRepo: repos.Auth,
		Redis:    rdb,
		Cache:    NewCacheService(rdb),
		Validate: validate,
	}

	return &Services{
		Users: NewUserService(deps),
		Decks: NewDeckService(deps),
		Auth:  NewAuthService(deps),
		Rdb:   rdb,
	}
}
