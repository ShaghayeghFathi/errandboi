package server

import (
	"context"
	"time"

	"errandboi/internal/config"
	"errandboi/internal/db/mongodb"
	"errandboi/internal/db/rdb"
	"errandboi/internal/http/handler"
	"errandboi/internal/publisher"
	"errandboi/internal/scheduler"
	"errandboi/internal/services/emq"
	"errandboi/internal/services/nats"
	"errandboi/internal/store/mongo"
	redisPK "errandboi/internal/store/redis"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func main(cfg config.Config, logger *zap.Logger) {
	ctx := context.Background()

	redisClient, err := rdb.New(ctx, cfg.Redis)
	if err != nil {
		logger.Fatal("redis initiation failed", zap.Error(err))
	}

	redis := redisPK.NewRedis(&rdb.Redis{Client: redisClient})

	mongodb, err := mongodb.New(cfg.Mongo)
	if err != nil {
		logger.Fatal("mongo initiation failed", zap.Error(err))
	}

	mongo := mongo.NewMongoDB(mongodb)

	app := fiber.New(fiber.Config{
		AppName: "errandboi",
	})

	handler.Handler{
		Redis:  redis,
		Mongo:  mongo,
		Logger: logger,
	}.Register(app)

	emqClient, err := emq.NewConnection(cfg.Emq)
	if err != nil {
		logger.Fatal("emq client initiation failed", zap.Error(err))
	}

	natsClient, err := nats.NewConnection(cfg.Nats, logger)
	if err != nil {
		logger.Fatal("nats client initiation failed", zap.Error(err))
	}

	// if err := natsClient.CreateStream(); err != nil {
	// 	logger.Fatal("stream creation failed", zap.Error(err))
	// }

	const ws = 10
	publisher := publisher.NewPublisher(redis, &emq.Mqtt{Client: emqClient}, natsClient, mongo, ws, logger)

	scheduler, err := scheduler.NewScheduler(publisher, logger)
	if err != nil {
		logger.Fatal("scheduler initiation failed", zap.Error(err))
	}

	scheduler.WorkInIntervals(time.Second)

	err = app.Listen(":3000")

	if err != nil {
		logger.Fatal("fiber initiation failed", zap.Error(err))
	}
}

func Register(root *cobra.Command, cfg config.Config, logger *zap.Logger) {
	root.AddCommand(
		&cobra.Command{
			Use:   "serve",
			Short: "Run server",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg, logger)
			},
		},
	)
}
