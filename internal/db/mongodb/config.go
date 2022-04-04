package mongodb

import "time"

type Config struct {
	Name    string        `koanf:"name"`
	URL     string        `koanf:"url"`
	Timeout time.Duration `koanf:"timeout"`
}
