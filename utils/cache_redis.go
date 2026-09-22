package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type FetchFunc[T any] func() (T, error)

func GetOrSetCache[T any](ctx context.Context, redisClient *redis.Client, key string, ttl time.Duration, fetch FetchFunc[T]) (T, error) {
	var result T

	cachedData, err := redisClient.Get(ctx, key).Result()
	if err == nil {
		if errJSON := json.Unmarshal([]byte(cachedData), &result); errJSON == nil {
			fmt.Println("Data dari Cache")
			return result, nil
		}
	}

	fmt.Println("Data dari Database")
	result, err = fetch()
	if err != nil {
		return result, err
	}

	dataJSON, _ := json.Marshal(result)
	redisClient.Set(ctx, key, dataJSON, ttl)

	return result, nil
}