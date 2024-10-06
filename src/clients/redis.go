package clients

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisClient struct {
	Client       *redis.Client
	RedisTimeout time.Duration
}

func NewRedisClient(address, password string, db int, timeout time.Duration) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	return &RedisClient{
		Client:       rdb,
		RedisTimeout: timeout,
	}
}

func (r *RedisClient) Set(key string, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := r.Client.Set(ctx, key, value, time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisClient) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	val, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}
