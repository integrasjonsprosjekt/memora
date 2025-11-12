package services

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	UserTTL     = 5 * time.Minute
	DeckTTL     = 5 * time.Minute
	CardTTL     = 10 * time.Minute
	DeckListTTL = 2 * time.Minute
	CardListTTL = 2 * time.Minute

	CacheOpTimeout = 5 * time.Second
)

type CacheService struct {
	rdb *redis.Client
}

func NewCacheService(rdb *redis.Client) *CacheService {
	return &CacheService{rdb: rdb}
}

func (c *CacheService) Get(ctx context.Context, key string, dest any) error {
	data, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), dest)
}

func (c *CacheService) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		slog.Error("failed to marshal cache value", "error", err)
		return err
	}
	err = c.rdb.Set(ctx, key, data, ttl).Err()
	if err != nil {
		slog.Error("failed to set cache value", "error", err)
	}
	return nil
}

func (c *CacheService) SetAsync(key string, value any, ttl time.Duration) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), CacheOpTimeout)
		defer cancel()
		err := c.Set(ctx, key, value, ttl)
		if err != nil {
			slog.Error("failed to set cache value asynchronously", "error", err)
		}
	}()
}

func (c *CacheService) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	err := c.rdb.Del(ctx, keys...).Err()
	if err != nil {
		slog.Error("failed to delete cache keys", "error", err)
	}
	return nil
}

func (c *CacheService) DeleteAsync(keys ...string) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), CacheOpTimeout)
		defer cancel()
		err := c.Delete(ctx, keys...)
		if err != nil {
			slog.Error("failed to delete cache keys asynchronously", "error", err)
		}
	}()
}

func (c *CacheService) DeletePattern(ctx context.Context, pattern string) error {
	iter := c.rdb.Scan(ctx, 0, pattern, 0).Iterator()

	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		slog.Error("failed to scan cache keys", "error", err)
		return err
	}

	return c.Delete(ctx, keys...)
}
