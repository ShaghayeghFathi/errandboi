package config

import (
	"errandboi/internal/db/mongodb"
	"errandboi/internal/db/rdb"
	"time"
)

func Default() Config{
	return Config{
	Redis: rdb.Config{
		Address:  "ADDRESS:6379",
		Password: "",
		DB:       0,
	},
	Mongo: mongodb.Config{
		Name: "errandboi",
		URL: "mongodb://ADDRESS:27017",
		Timeout: 10*time.Second,
	},
   }	
}