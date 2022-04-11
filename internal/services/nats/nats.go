package nats

import (
	"fmt"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Nats struct {
	Connection *nats.Conn
	JSCtx      nats.JetStreamContext
	Logger     *zap.Logger
}

func NewConnection(cfg Config, logger *zap.Logger) (*Nats, error) {
	nc, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("nats connection failed %w", err)
	}

	nc.SetDisconnectErrHandler(func(c *nats.Conn, err error) {
		logger.Fatal("nats disconnected", zap.Error(err))
	})

	nc.SetReconnectHandler(func(c *nats.Conn) {
		logger.Warn("nats reconnected")
	})

	jsc, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("could not get jet stream context %w", err)
	}

	return &Nats{
		Connection: nc,
		JSCtx:      jsc,
		Logger:     logger,
	}, nil
}

func (n *Nats) CreateStream() error {
	stream, _ := n.JSCtx.StreamInfo(ChannelName)

	if stream == nil {
		in, err2 := n.JSCtx.AddStream(&nats.StreamConfig{
			Name:     ChannelName,
			Subjects: []string{SubjectName},
			MaxAge:   0,
			Storage:  nats.FileStorage,
		})
		if err2 != nil {
			return fmt.Errorf("cannot create stream %w", err2)
		}

		stream = in
	}

	n.Logger.Info("events stream", zap.Any("stream", stream))

	return nil
}
