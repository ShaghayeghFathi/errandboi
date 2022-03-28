package nats

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

type Nats struct {
	Connection *nats.Conn
	JContext   nats.JetStreamContext
}

func NewConnection(cfg Config) (*Nats, error) {
	fmt.Printf(cfg.URL)
	nc, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("nats connection failed %w", err)
	}

	nc.SetDisconnectErrHandler(func(c *nats.Conn, err error) {
		log.Fatal("nats disconnected", err)
	})

	nc.SetReconnectHandler(func(c *nats.Conn) {
		log.Println("nats reconnected")
	})

	jsc, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("could not get jet stream context %w", err)
	}
	return &Nats{
		Connection: nc,
		JContext:   jsc,
	}, nil

}
