package services

import (
	"context"
	"memora/internal/firebase"
	"memora/internal/models"

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

func (s *UserService) GetUser(ctx context.Context, id string) (models.User, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *UserService) RegisterNewUser(ctx context.Context, user models.CreateUser) (string, error) {
	id, err := s.repo.AddUser(ctx, user)
	if err != nil {
		return "", err
	}

	return id, nil
}
