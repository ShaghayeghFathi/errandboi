package server

import (
	"context"
	"errandboi/internal/config"
	"errandboi/internal/db/mongodb"
	"errandboi/internal/db/rdb"
	"errandboi/internal/http/handler"
	"errandboi/internal/store/mongo"
	redisPK "errandboi/internal/store/redis"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
)

func main(cfg config.Config){
	
	println("ran serve command")
	ctx := context.Background()

	redisClient, err := rdb.New(ctx, cfg.Redis)
	if err != nil {
		fmt.Println("redis initiation failed")

	}
	redisdb := rdb.Redis{Client: redisClient}
	mongodb, err := mongodb.New(cfg.Mongo)
	if err != nil {
		log.Fatal("mongo initiation failed")
	}
	app := fiber.New(fiber.Config{
		AppName: "errandboi",
	})

	handler.Handler{
		Redis : redisPK.NewRedis(&redisdb),
		Mongo: mongo.NewMongoDB(mongodb),
	}.Register(app)

	if err := app.Listen(":3000"); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("fiber initiation failed")
	}



}

func Register(root *cobra.Command, cfg config.Config) {
	root.AddCommand(
		&cobra.Command{
			Use:   "serve",
			Short: "Run server",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg)
			},
		},
	)
}