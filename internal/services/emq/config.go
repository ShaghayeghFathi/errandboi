package emq

type Config struct {
	Broker   string `koanf:"broker"`
	Port     int    `koanf:"port"`
	ClientID string `koanf:"clientid"`
	Username string `koanf:"username"`
	Password string `koanf:"password"`
}
