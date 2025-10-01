package services

import (
	"context"
	"memora/internal/errors"
	"memora/internal/firebase"
	"memora/internal/models"
	"memora/internal/utils"

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

func (s *UserService) UpdateUser(
	ctx context.Context,
	updateStruct models.PatchUser,
	id string,
) error {
	if err := s.validate.Struct(updateStruct); err != nil {
		return errors.ErrInvalidUser
	}

	update, err := utils.StructToUpdate(updateStruct)
	if err != nil {
		return errors.ErrInvalidUser
	}

	err = s.repo.UpdateUser(ctx, update, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.DeleteUser(ctx, id)
}
