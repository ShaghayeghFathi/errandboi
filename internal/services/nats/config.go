package nats

type Config struct {
	URL string `koanf:"url"`
}

const (
	ChannelName = "events"
)
