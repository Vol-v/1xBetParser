package cache

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // or your Redis server address
		Password: "",               // no password set
		DB:       0,                // default DB
	})
}

func StoreInRedis(rdb *redis.Client, key string, data map[string]string) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return rdb.Set(ctx, key, jsonData, 0).Err()
}

func RetrieveFromRedis(rdb *redis.Client, key string) (map[string]string, error) {
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var data map[string]string
	err = json.Unmarshal([]byte(val), &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
