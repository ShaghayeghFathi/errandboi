package rdb

type Config struct {
	Address  string `koanf:"address"`
	Password string `koanf:"password"`
	DB       int    `koanf:"db"`
}
