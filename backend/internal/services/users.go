package services

import (
	"context"
	"memora/internal/errors"
	"memora/internal/firebase"
	"memora/internal/models"
	"memora/internal/utils"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
)

// Default filter for all fields, used when updating a user
var defaultFilterUsers = "email,name"

// UserService provides methods for managing users.
type UserService struct {
	repo     firebase.UserRepository
	rdb      *redis.Client
	validate *validator.Validate
}

// NewUserService creates a new instance of UserService.
func NewUserService(
	deps *ServiceDeps,
) *UserService {
	return &UserService{
		repo:     deps.UserRepo,
		rdb:      deps.Redis,
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
	cachedUser, err := utils.GetDataFromRedis[models.User](cacheKey, s.rdb, ctx)
	if err == nil {
		return cachedUser, nil
	}

	user, err := s.repo.GetUser(ctx, id, filterParsed)
	if err != nil {
		return models.User{}, err
	}

	// Store the user in the cache for future requests
	utils.SetDataToRedis(cacheKey, user, s.rdb, ctx, 5*time.Minute)

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
	cachedDecks, err := utils.GetDataFromRedis[models.UserDecks](cacheKey, s.rdb, ctx)
	if err == nil {
		return cachedDecks, nil
	}

	decks, err := s.repo.GetDecks(ctx, id, filterParsed)
	if err != nil {
		return models.UserDecks{}, err
	}

	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		utils.SetDataToRedis(cacheKey, decks, s.rdb, bgCtx, 5*time.Minute)
	}()

	// Return the list of decks
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

	cacheKey := utils.UserKey(id)

	utils.DeleteDataFromRedis(cacheKey, s.rdb, ctx)

	// Perform the update in the repository
	err = s.repo.UpdateUser(ctx, update, id)
	if err != nil {
		return models.User{}, err
	}

	return s.GetUser(ctx, id, defaultFilterUsers)
}

// DeleteUser removes a user by their ID.
// Returns an error if the operation fails.
func (s *UserService) DeleteUser(
	ctx context.Context,
	id string,
) error {
	utils.DeleteDataFromRedis(utils.UserKey(id), s.rdb, ctx)

	return s.repo.DeleteUser(ctx, id)
}
