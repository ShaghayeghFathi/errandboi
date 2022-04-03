package logger

type Config struct {
	Level  string `koanf:"level"`
	Syslog `koanf:"syslog"`
}

type Syslog struct {
	Enabled bool   `koanf:"enabled"`
	Network string `koanf:"network"`
	Address string `koanf:"address"`
	Tag     string `koanf:"tag"`
}
