package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func New(cfg Config) (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(cfg.URL)

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("new mongo client error: %w", err)
	}

	{
		ctx, done := context.WithTimeout(context.Background(), cfg.Timeout)
		defer done()

		if err := client.Connect(ctx); err != nil {
			return nil, fmt.Errorf("mongo connection error: %w", err)
		}
	}
	{
		ctx, done := context.WithTimeout(context.Background(), cfg.Timeout)
		defer done()

		if err := client.Ping(ctx, readpref.Primary()); err != nil {
			return nil, fmt.Errorf("mongo ping error: %w", err)
		}
	}

	return client.Database(cfg.Name), nil
}
