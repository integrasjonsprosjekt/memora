package utils

import (
	"context"
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

func SetDataToRedis[T any](key string, data T, rdb *redis.Client, ctx context.Context, ttl time.Duration) error {
	return rdb.Set(ctx, key, data, ttl).Err()
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
