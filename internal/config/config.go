package config

import (
	"errandboi/internal/db/mongodb"
	"errandboi/internal/db/rdb"
	"log"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/structs"
)

type Config struct{
	Redis rdb.Config `koanf:"redis"`
	Mongo mongodb.Config `koanf:"mongo"`
}

func New() Config {
	k := koanf.New(".")

	if err := k.Load(structs.Provider(Default(), "koanf"), nil); err != nil {
		log.Fatalf("error loading default config: %v", err)
	}
	var cfg Config
	err := k.Unmarshal("", &cfg)
	if err != nil {
		log.Fatalf("erro unmarshaling: %v", err)
	}

	return cfg
}