package publisher

import (
	"context"
	"errandboi/internal/services/emq"
	"errandboi/internal/services/nats"
	"errandboi/internal/store/mongo"
	redisPK "errandboi/internal/store/redis"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gammazero/workerpool"
)

type Publisher struct {
	Redis      *redisPK.RedisDB
	Mongo      *mongo.MongoDB
	Events     []Event
	Mqtt       *emq.Mqtt
	Nats       *nats.Nats
	Wp         *workerpool.WorkerPool
	WorkerSize int
}

type Event struct {
	ID      string   `json:"id"`
	Topic   string   `json:"topic"`
	Payload string   `json:"payload"`
	Type    []string `json:"type"`
}

var EventRedisFields []string = []string{"topic", "payload", "type"}

func NewPublisher(r *redisPK.RedisDB, client *emq.Mqtt, natsCl *nats.Nats, m *mongo.MongoDB, size int) *Publisher {
	return &Publisher{Redis: r, Mongo: m, Mqtt: client, Nats: natsCl, Wp: workerpool.New(size), WorkerSize: size}
}

func (pb *Publisher) GetEvents() {
	var ctx = context.Background()
	start := float64(time.Now().Unix())
	events, err := pb.Redis.ZGetRange(ctx, "events", start, start+1)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(events); i++ {
		eventId := events[i].Member.(string)
		var field []string
		for j := 0; j < 3; j++ {
			field[i], err = pb.Redis.Get(ctx, EventRedisFields[i]+"_"+eventId)
			if err != nil {
				log.Fatal("Could not retrieve event info from redis")
			}
		}

		types := strings.Split(field[2], "_")
		pb.Events = append(pb.Events, Event{ID: eventId, Topic: field[0], Payload: field[1], Type: types})
	}
}

func (pb *Publisher) RemoveEvent(event Event) error {
	var index = -1
	for i := range pb.Events {
		if pb.Events[i].ID == event.ID {
			index = i
		}
	}
	if index == -1 {
		return errors.New("event to be deleted was not found in publisher events")
	}
	pb.Events[index], pb.Events[len(pb.Events)-1] = pb.Events[len(pb.Events)-1], pb.Events[index]
	pb.Events = pb.Events[:len(pb.Events)-1]

	pb.deleteEventRedis(event)

	// pb.Mongo.updateEventStatus(event.ID)

	return nil
}

func (pb *Publisher) Cancel() error {
	pb.Wp.Stop()
	if !pb.Wp.Stopped() {
		return errors.New("could not stop publisher")
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
			for i := 0; i < len(event.Type); i++ {
				pb.publishEvent(event.Type[i], event)
			}
		})
	}
	wg.Wait()
}

func (pb *Publisher) publishEvent(t string, event Event) {
	if t == "emqx" {
		pb.publishEventEMQ(event)
	} else {
		pb.publishEventNats(event)
	}
}

func (pb *Publisher) publishEventEMQ(event Event) {
	if token := pb.Mqtt.Client.Subscribe(event.Topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
	token := pb.Mqtt.Client.Publish(event.Topic, 0, false, event.Payload)
	token.Wait()
}

func (pb *Publisher) publishEventNats(event Event) {

}

func (pb *Publisher) deleteEventRedis(e Event) {
	ctx := context.Background()
	for i := 0; i < 3; i++ {
		err := pb.Redis.Del(ctx, EventRedisFields[i]+"_"+e.ID)
		if err != nil {
			log.Fatal("Could delete event info from redis")
		}
	}
}
