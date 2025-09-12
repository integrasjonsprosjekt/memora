package services

import (
	"context"
	"memora/internal/models"
)

type UserService struct {
	repo models.UserRepository
}

func NewUserService(repo models.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterNewUser(ctx context.Context, user models.CreateUser) (string, error) {
	id, err := s.repo.AddUser(ctx, user)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *UserService) GetUser(ctx context.Context, id string) (models.User, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
