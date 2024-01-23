package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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

func StoreInRedis(rdb *redis.Client, key string, data map[string]string, duration time.Duration) error {
	/*I am storing data with set because when I retrieve I need to get all {bet_name : coefficient} pairs*/
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	fmt.Printf("storing in redis...\n")
	return rdb.Set(ctx, key, jsonData, duration).Err()
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
