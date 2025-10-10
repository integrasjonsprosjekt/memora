package services

import (
	"context"
	"memora/internal/errors"
	"memora/internal/firebase"
	"memora/internal/models"
	"memora/internal/utils"

	"github.com/go-playground/validator/v10"
)

// UserService provides methods for managing users.
type UserService struct {
	repo     firebase.UserRepository
	validate *validator.Validate
}

// NewUserService creates a new instance of UserService.
func NewUserService(
	repo firebase.UserRepository,
	validate *validator.Validate,
) *UserService {
	return &UserService{
		repo:     repo,
		validate: validate,
	}
}

// GetUser retrieves a user by their ID.
// Returns the user or an error if the operation fails.
func (s *UserService) GetUser(
	ctx context.Context,
	id string,
) (models.User, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

// GetDecksOwned retrieves all decks owned by a user.
// Returns a list of decks or an error if the operation fails.
func (s *UserService) GetDecksOwned(
	ctx context.Context,
	id string,
) ([]models.DisplayDeck, error) {
	decks, err := s.repo.GetDecksOwned(ctx, id)

	if decks == nil {
		return []models.DisplayDeck{}, err
	}

	// Return the list of decks
	return decks, err
}

// GetDecksShared retrieves all decks shared with a user.
// Returns a list of decks or an error if the operation fails.
func (s *UserService) GetDecksShared(
	ctx context.Context,
	id string,
) ([]models.DisplayDeck, error) {
	decks, err := s.repo.GetDecksShared(ctx, id)

	if decks == nil {
		return []models.DisplayDeck{}, err
	}

	return decks, err
}

// RegisterNewUser creates a new user from the provided data.
// Returns the new user's ID or an error if the operation fails.
func (s *UserService) RegisterNewUser(
	ctx context.Context,
	user models.CreateUser,
) (string, error) {
	id, err := s.repo.AddUser(ctx, user)
	if err != nil {
		return "", err
	}

	return id, nil
}

// UpdateUser updates fields of an existing user.
// Validates the input and returns an error if the operation fails.
func (s *UserService) UpdateUser(
	ctx context.Context,
	updateStruct models.PatchUser,
	id string,
) error {
	// Validate the input struct
	if err := s.validate.Struct(updateStruct); err != nil {
		return errors.ErrInvalidUser
	}

	// Convert the struct to Firestore update format
	update, err := utils.StructToUpdate(updateStruct)
	if err != nil {
		return errors.ErrInvalidUser
	}

	// Perform the update in the repository
	err = s.repo.UpdateUser(ctx, update, id)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUser removes a user by their ID.
// Returns an error if the operation fails.
func (s *UserService) DeleteUser(
	ctx context.Context,
	id string,
) error {
	return s.repo.DeleteUser(ctx, id)
}
