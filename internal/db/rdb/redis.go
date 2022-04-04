package rdb

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

type Redis struct {
	Client *redis.Client
}

func New(ctx context.Context, cfg Config) (*redis.Client, error) {
	const t = 10
	timeout, cancel := context.WithTimeout(ctx, t*time.Second)

	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(timeout).Err(); err != nil {
		return nil, fmt.Errorf("error connecting to redis: %w", err)
	}

	return client, nil
}
