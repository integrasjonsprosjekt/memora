package services

import (
	"memora/internal/firebase"

	"github.com/go-playground/validator/v10"
)

type Services struct {
	Users *UserService
}

func NewServices(repos *firebase.Repositories, validate *validator.Validate) *Services {
	return &Services{
		Users: NewUserService(repos.User, validate),
	}
}
