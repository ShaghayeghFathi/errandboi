package publisher

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ShaghayeghFathi/errandboi/internal/services/emq"
	natsp "github.com/ShaghayeghFathi/errandboi/internal/services/nats"
	"github.com/ShaghayeghFathi/errandboi/internal/store/mongo"
	redisp "github.com/ShaghayeghFathi/errandboi/internal/store/redis"

	"github.com/gammazero/workerpool"
	"go.uber.org/zap"
)

type Publisher struct {
	Redis      *redisp.RedisDB
	Mongo      *mongo.DB
	Events     []Event
	Mqtt       *emq.Mqtt
	Nats       *natsp.Nats
	Wp         *workerpool.WorkerPool
	WorkerSize int
	logger     *zap.Logger
}

type Event struct {
	ID      string   `json:"id"`
	Topic   string   `json:"topic"`
	Payload string   `json:"payload"`
	Type    []string `json:"type"`
}

const setName = "events"

func NewPublisher(r *redisp.RedisDB, client *emq.Mqtt, natsCl *natsp.Nats,
	m *mongo.DB, size int, logger *zap.Logger,
) *Publisher {
	return &Publisher{
		Redis: r, Mongo: m, Mqtt: client, Nats: natsCl,
		Wp: workerpool.New(size), WorkerSize: size, logger: logger,
	}
}

func (pb *Publisher) GetEvents() {
	eventRedisFields := []string{"topic", "payload", "type"}

	pb.Events = make([]Event, 0)

	ctx := context.Background()

	start := float64(time.Now().Unix())

	events, err := pb.Redis.ZGetRange(ctx, setName, start, start+1)
	if err != nil {
		pb.logger.Info("could not retrieve event", zap.Error(err))
	}

	for i := 0; i < len(events); i++ {
		eventID, _ := events[i].Member.(string)

		var field []string

		for j := 0; j < 3; j++ {
			tmp, err := pb.Redis.Get(ctx, eventRedisFields[j]+"_"+eventID)
			if err != nil {
				pb.logger.Warn("could not retrirve event fields", zap.Error(err))
			}

			field = append(field, tmp)
		}

		types := strings.Split(field[2], "_")

		pb.Events = append(pb.Events, Event{ID: eventID, Topic: field[0], Payload: field[1], Type: types})
		pb.deleteEventRedis(eventID)
	}
}

func (pb *Publisher) Cancel() {
	pb.Wp.Stop()

	if !pb.Wp.Stopped() {
		pb.logger.Warn("publisher not stopped")
	}
}

func (pb *Publisher) Work() {
	var wg sync.WaitGroup

	for idx := range pb.Events {
		event := pb.Events[idx]

		wg.Add(1)
		pb.Wp.Submit(func() {
			defer wg.Done()
			pb.publishEvent(event)
		})
		// update only events status
		pb.Mongo.UpdateEventStatus(context.Background(), event.ID)
	}

	wg.Wait()
}

func (pb *Publisher) publishEvent(event Event) {
	for i := 0; i < len(event.Type); i++ {
		if event.Type[i] == "emqx" {
			go pb.publishEventEMQ(event)
		} else if event.Type[i] == "nats" {
			go pb.publishEventNats(event)
		}
	}
}

func (pb *Publisher) publishEventEMQ(event Event) {
	if token := pb.Mqtt.Client.Subscribe(event.Topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}

	token := pb.Mqtt.Client.Publish(event.Topic, 0, false, event.Payload)
	pb.logger.Info("message published to emq: ", zap.String("payload", event.Payload))
	token.Wait()
}

func (pb *Publisher) publishEventNats(event Event) {
	t := natsp.ChannelName + "." + event.Topic
	if _, err := pb.Nats.JSCtx.Publish(t, []byte(event.Payload)); err != nil {
		pb.logger.Error("failed to publish event", zap.String("payload", event.Payload),
			zap.String("topic", event.Topic), zap.Error(err))
	}

	pb.logger.Info("message published to nats", zap.String("payload", event.Payload))
}

func (pb *Publisher) deleteEventRedis(id string) {
	eventRedisFields := []string{"topic", "payload", "type"}

	ctx := context.Background()

	err, _ := pb.Redis.ZRem(ctx, "events", id)
	if err == 0 {
		pb.logger.Warn("failed to delete event from redis")
	}

	for i := 0; i < 3; i++ {
		err := pb.Redis.Del(ctx, eventRedisFields[i]+"_"+id)
		if err != nil {
			pb.logger.Error("failed to delete event from redis", zap.Error(err))
		}
	}
}
