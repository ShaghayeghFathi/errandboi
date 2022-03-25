package config

import (
	"errandboi/internal/db/mongodb"
	"errandboi/internal/db/rdb"
	"time"
)

func Default() Config{
	return Config{
	Redis: rdb.Config{
		Address:  "localhost:6379",
		Password: "",
		DB:       0,
	},
	Mongo: mongodb.Config{
		Name: "errandboi",
		URL: "mongodb://localhost:27017",
		Timeout: 10*time.Second,
	},
   }	
}