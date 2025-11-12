package services

import (
	"context"
	"memora/internal/errors"
	"memora/internal/firebase"
	"memora/internal/models"
	"memora/internal/utils"

	"github.com/go-playground/validator/v10"
)

// Default filter for all fields, used when updating a user
var defaultFilterUsers = "email,name"

// UserService provides methods for managing users.
type UserService struct {
	repo     firebase.UserRepository
	cache    *CacheService
	validate *validator.Validate
}

// NewUserService creates a new instance of UserService.
func NewUserService(
	deps *ServiceDeps,
) *UserService {
	return &UserService{
		repo:     deps.UserRepo,
		cache:    deps.Cache,
		validate: deps.Validate,
	}
}

// GetUser retrieves a user by their ID.
// Filter specifies which fields to return.
// Returns the user or an error if the operation fails.
func (s *UserService) GetUser(
	ctx context.Context,
	id, filter string,
) (models.User, error) {

	filterParsed, err := utils.ParseFilter(filter)
	if err != nil {
		return models.User{}, err
	}

	cacheKey := utils.UserKey(id)

	// Try to get the user from the cache
	var user models.User
	err = s.cache.Get(ctx, cacheKey, &user)
	if err == nil {
		return user, nil
	}

	user, err = s.repo.GetUser(ctx, id, filterParsed)
	if err != nil {
		return models.User{}, err
	}

	// Store the user in the cache for future requests
	s.cache.SetAsync(cacheKey, user, UserTTL)

	return user, nil
}

// GetDecks retrieves all decks associated with a user.
// Filter specifies which fields to return.
// Returns a list of decks or an error if the operation fails.
func (s *UserService) GetDecks(
	ctx context.Context,
	id, filter, email string,
) (models.UserDecks, error) {
	filterParsed, err := utils.ParseFilter(filter)
	if err != nil {
		return models.UserDecks{}, err
	}

	cacheKey := utils.UserEmailDecksKey(email)
	var decks models.UserDecks
	err = s.cache.Get(ctx, cacheKey, &decks)
	if err == nil {
		return decks, nil
	}

	decks, err = s.repo.GetDecks(ctx, id, filterParsed)
	if err != nil {
		return models.UserDecks{}, err
	}

	s.cache.SetAsync(cacheKey, decks, UserTTL)

	return decks, nil
}

// RegisterNewUser creates a new user from the provided data.
// Returns the new user's ID or an error if the operation fails.
func (s *UserService) RegisterNewUser(
	ctx context.Context,
	user models.CreateUser,
	id string,
) error {
	// Validate the input struct
	if err := s.validate.Struct(user); err != nil {
		return errors.ErrInvalidUser
	}

	return s.repo.AddUser(ctx, user, id)
}

// UpdateUser updates fields of an existing user.
// Validates the input and returns an error if the operation fails.
func (s *UserService) UpdateUser(
	ctx context.Context,
	updateStruct models.PatchUser,
	id string,
) (models.User, error) {
	// Validate the input struct
	if err := s.validate.Struct(updateStruct); err != nil {
		return models.User{}, errors.ErrInvalidUser
	}

	// Convert the struct to Firestore update format
	update, err := utils.StructToUpdate(updateStruct)
	if err != nil {
		return models.User{}, errors.ErrInvalidUser
	}

	// Perform the update in the repository
	err = s.repo.UpdateUser(ctx, update, id)
	if err != nil {
		return models.User{}, err
	}

	s.cache.Delete(ctx, utils.UserKey(id))

	return s.GetUser(ctx, id, defaultFilterUsers)
}

// DeleteUser removes a user by their ID.
// Returns an error if the operation fails.
func (s *UserService) DeleteUser(
	ctx context.Context,
	id string,
) error {
	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	s.cache.Delete(ctx, utils.UserKey(id))

	return nil
}
