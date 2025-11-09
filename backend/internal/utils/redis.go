package utils

import (
	"context"
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
	err := rdb.Get(ctx, key).Scan(&result)
	return result, err
}

func SetDataToRedis[T any](key string, data T, rdb *redis.Client, ctx context.Context, ttl time.Duration) {
	err := rdb.Set(ctx, key, data, ttl).Err()
	if err != nil {
		slog.Error("Error setting data to Redis", slog.Any("err", err))
	}
}

func DeleteDataFromRedis(key string, rdb *redis.Client, ctx context.Context) error {
	return rdb.Del(ctx, key).Err()
}

func UserKey(userID string) string {
	return UserKeyPrefix + ":" + userID
}

func DeckKey(deckID string) string {
	return DeckKeyPrefix + ":" + deckID
}

func CardKey(cardID string) string {
	return CardKeyPrefix + ":" + cardID
}
