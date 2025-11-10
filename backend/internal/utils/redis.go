package utils

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	UserKeyPrefix = "user"
	DeckKeyPrefix = "deck"
	CardKeyPrefix = "card"
)

func GetDataFromRedis[T any](key string, rdb *redis.Client, ctx context.Context) (T, error) {
	var result T

	data, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(data), &result)
	return result, err
}

func SetDataToRedis[T any](
	key string,
	data T,
	rdb *redis.Client,
	ctx context.Context,
	ttl time.Duration,
) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		slog.Error("Error marshalling data to JSON", slog.Any("err", err))
		return
	}

	err = rdb.Set(ctx, key, jsonData, ttl).Err()
	if err != nil {
		slog.Error("Error setting data to Redis", slog.Any("err", err))
	}
}

func DeleteDataFromRedis(key string, rdb *redis.Client, ctx context.Context) {
	err := rdb.Del(ctx, key).Err()
	if err != nil {
		slog.Error("Error deleting data from Redis", slog.Any("err", err))
	}
}

func UserKey(userID string) string {
	return UserKeyPrefix + ":" + userID
}

func DeckKey(deckID string) string {
	return DeckKeyPrefix + ":" + deckID
}

func DeckCardKey(deckID, cardID string) string {
	return DeckKeyPrefix + ":" + deckID + ":" + CardKeyPrefix + ":" + cardID
}

func UserEmailDecksKey(email string) string {
	return "user:email:" + email + ":decks"
}

func UserKeyRateLimit(userID string) string {
	return "rate_limit:user:" + userID
}
