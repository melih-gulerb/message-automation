package clients

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisClient struct {
	RedisAddress  string
	RedisPassword string
	RedisDB       int
	RedisTimeout  time.Duration
}

func NewRedisClient(address, password string, db int, timeout time.Duration) *RedisClient {
	return &RedisClient{
		RedisAddress:  address,
		RedisPassword: password,
		RedisDB:       db,
		RedisTimeout:  timeout,
	}
}

func (r *RedisClient) Set(rdb *redis.Client, key string, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := rdb.Set(ctx, key, value, time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisClient) Get(rdb *redis.Client, key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}
