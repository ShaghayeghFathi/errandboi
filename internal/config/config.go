package config

import (
	"log"

	"github.com/ShaghayeghFathi/errandboi/internal/db/mongodb"
	"github.com/ShaghayeghFathi/errandboi/internal/db/rdb"
	"github.com/ShaghayeghFathi/errandboi/internal/logger"
	"github.com/ShaghayeghFathi/errandboi/internal/services/emq"
	"github.com/ShaghayeghFathi/errandboi/internal/services/nats"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/structs"
)

type Config struct {
	Logger logger.Config  `koanf:"logger"`
	Redis  rdb.Config     `koanf:"redis"`
	Mongo  mongodb.Config `koanf:"mongo"`
	Emq    emq.Config     `koanf:"emq"`
	Nats   nats.Config    `koanf:"nats"`
}

func New() Config {
	k := koanf.New(".")

	if err := k.Load(structs.Provider(Default(), "koanf"), nil); err != nil {
		log.Fatalf("error loading default config: %v", err)
	}

	var cfg Config

	if err := k.Unmarshal("", &cfg); err != nil {
		log.Fatalf("erro unmarshaling: %v", err)
	}

	return cfg
}
