package services

import (
	"memora/internal/firebase"

	"github.com/go-playground/validator/v10"
)

type UserService struct {
	repo     firebase.UserRepository
	validate *validator.Validate
}

func NewUserService(repo firebase.UserRepository, validate *validator.Validate) *UserService {
	return &UserService{repo: repo, validate: validate}
}
