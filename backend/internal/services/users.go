package services

import (
	"context"
	customerror "memora/internal/customError"
	"memora/internal/errors"
	"memora/internal/firebase"
	"memora/internal/models"
	"memora/internal/utils"

	"github.com/go-playground/validator/v10"
)

type UserService struct {
	repo     firebase.UserRepository
	validate *validator.Validate
}

func NewUserService(repo firebase.UserRepository, validate *validator.Validate) *UserService {
	return &UserService{repo: repo, validate: validate}
}

func (s *UserService) GetUser(ctx context.Context, id string) (models.User, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
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

func (s *UserService) UpdateUser(ctx context.Context, update map[string]interface{}, id string) error {
	if !validatePatch(update) {
		return customerror.ErrInvalidUser
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
