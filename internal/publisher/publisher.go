package publisher

import (
	"context"
	"errandboi/internal/services/emq"
	natsPK "errandboi/internal/services/nats"
	"errandboi/internal/store/mongo"
	redisPK "errandboi/internal/store/redis"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gammazero/workerpool"
	"go.uber.org/zap"
)

type Publisher struct {
	Redis      *redisPK.RedisDB
	Mongo      *mongo.MongoDB
	Events     []Event
	Mqtt       *emq.Mqtt
	Nats       *natsPK.Nats
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

var EventRedisFields []string = []string{"topic", "payload", "type"}

const setName = "events"

func NewPublisher(r *redisPK.RedisDB, client *emq.Mqtt, natsCl *natsPK.Nats, m *mongo.MongoDB, size int, logger *zap.Logger) *Publisher {
	return &Publisher{Redis: r, Mongo: m, Mqtt: client, Nats: natsCl, Wp: workerpool.New(size), WorkerSize: size, logger: logger}
}

func (pb *Publisher) GetEvents() {
	pb.Events = make([]Event, 0)

	var ctx = context.Background()
	start := float64(time.Now().Unix())

	events, err := pb.Redis.ZGetRange(ctx, setName, start, start+1)
	if err != nil {
		pb.logger.Info("could not retrieve event", zap.Error(err))
	}

	for i := 0; i < len(events); i++ {
		eventId := events[i].Member.(string)
		var field []string
		for j := 0; j < 3; j++ {
			tmp, err := pb.Redis.Get(ctx, EventRedisFields[j]+"_"+eventId)
			if err != nil {
				pb.logger.Warn("could not retrirve event fields", zap.Error(err))
			}
			field = append(field, tmp)
		}

		types := strings.Split(field[2], "_")
		pb.Events = append(pb.Events, Event{ID: eventId, Topic: field[0], Payload: field[1], Type: types})
		pb.deleteEventRedis(eventId)
	}
}

func (pb *Publisher) Cancel() error {

	pb.Wp.Stop()
	if !pb.Wp.Stopped() {
		return fmt.Errorf("could not stop publisher")
	}
	return nil
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

	t := natsPK.ChannelName + "." + event.Topic

	if _, err := pb.Nats.JSCtx.Publish(t, []byte(event.Payload)); err != nil {
		pb.logger.Error("failed to publish event", zap.String("payload", event.Payload), zap.String("topic", event.Topic), zap.Error(err))
	}
	pb.logger.Info("message published to nats", zap.String("payload", event.Payload))
}

func (pb *Publisher) deleteEventRedis(id string) {

	ctx := context.Background()

	for i := 0; i < 3; i++ {
		err := pb.Redis.Del(ctx, EventRedisFields[i]+"_"+id)
		if err != nil {
			pb.logger.Error("failed to delete event from redis", zap.Error(err))
		}
	}
}

// func (pb *Publisher) RemoveEvent(event Event) error {
// 	var index = -1
// 	for i := range pb.Events {
// 		if pb.Events[i].ID == event.ID {
// 			index = i
// 		}
// 	}
// 	if index == -1 {
// 		return fmt.Errorf("event to be deleted was not found in publisher events")
// 	}
// 	pb.Events[index], pb.Events[len(pb.Events)-1] = pb.Events[len(pb.Events)-1], pb.Events[index]
// 	pb.Events = pb.Events[:len(pb.Events)-1]

// 	fmt.Println(pb.Events)

// 	return nil
// }
