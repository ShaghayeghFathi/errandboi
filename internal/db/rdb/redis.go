package rdb

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/gommon/log"
	"golang.org/x/net/context"
)

type Redis struct{
	Client *redis.Client
}

func New(ctx context.Context,cfg Config)(*redis.Client,error){
	timeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	log.Info("Connecting to redis...")
	if err := client.Ping(timeout).Err(); err != nil {
		return nil, fmt.Errorf("error connecting to redis: %v", err)
	}
	return client,nil
}


