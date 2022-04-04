package config

import (
	"time"

	"errandboi/internal/db/mongodb"
	"errandboi/internal/db/rdb"
	"errandboi/internal/logger"
	"errandboi/internal/services/emq"
	"errandboi/internal/services/nats"
)

func Default() Config {
	const t = 10

	const p = 1883

	return Config{
		Redis: rdb.Config{
			Address:  "localhost:6379",
			Password: "",
			DB:       0,
		},
		Mongo: mongodb.Config{
			Name:    "errandboi",
			URL:     "mongodb://localhost:27017",
			Timeout: t * time.Second,
		},
		Emq: emq.Config{
			Broker:   "localhost",
			Port:     p,
			ClientID: "go_mqtt_client",
			Username: "emqx",
			Password: "public",
		},
		Nats: nats.Config{
			URL: "nats://0.0.0.0:4222",
		},
		Logger: logger.Config{
			Level: "debug",
			Syslog: logger.Syslog{
				Enabled: false,
				Network: "",
				Address: "",
				Tag:     "",
			},
		},
	}
}
