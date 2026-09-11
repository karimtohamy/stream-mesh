package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func Get[T any](ctx context.Context, rdn *redis.Client, key string) (*T, error) {
	val, err := rdn.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var result T
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func Set[T any](ctx context.Context, rdn *redis.Client, key string, value T, ttl time.Duration) error {
	val, err := json.Marshal(value)
	fmt.Sprintf("Cache Set")
	if err != nil {
		return err
	}
	return rdn.Set(ctx, key, string(val), ttl).Err()
}
