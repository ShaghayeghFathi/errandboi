package publisher

import (
	"context"
	"errandboi/internal/services/emq"
	"errandboi/internal/store/mongo"
	redisPK "errandboi/internal/store/redis"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gammazero/workerpool"
)

type Publisher struct{
	Redis *redisPK.RedisDB
	Mongo *mongo.MongoDB
	Events []Event
	Mqtt *emq.Mqtt
	Wp         *workerpool.WorkerPool
	WorkerSize int
}

type Event struct{
	Topic string `json:"topic"`
	Payload string `json:"payload"`
	Type []string `json:"type"` 
}

func NewPublisher(r *redisPK.RedisDB, client *emq.Mqtt, m *mongo.MongoDB, size int) *Publisher{
	return &Publisher{Redis: r,Mongo: m, Mqtt : client,  Wp: workerpool.New(size), WorkerSize: size}
}

func(pb *Publisher) GetEvents(){
	var ctx = context.Background()
	start := float64(time.Now().Unix())
	events, err := pb.Redis.ZGetRange(ctx, "events", start, start+1)
	if err!=nil{
		log.Fatal(err)
	}
	for i := 0; i < len(events); i++ {
		eventId := events[i].Member.(string)
		topic , err := pb.Redis.Get(ctx, "topic" + "_" + eventId )
		if err !=nil{
			log.Fatal("Could not retrieve event topic from redis")
		}

		payload , err := pb.Redis.Get(ctx, "payload" + "_" + eventId )
		if err !=nil{
			log.Fatal("Could not retrieve event payload from redis")
		}

		types , err := pb.Redis.Get(ctx, "type" + "_" + eventId )
		if err !=nil{
			log.Fatal("Could not retrieve event type from redis")
		}
		typeSl := strings.Split(types, "_")
		// s:= strings.Split(eventId, "_")
		// event, err := pb.Mongo.GetEvent(context.Background(), eventId, s[0])
		// if err!= nil{
		// 	log.Fatal(err)
		// }
		pb.Events = append(pb.Events, Event{Topic: topic, Payload: payload , Type: typeSl})
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
	}
	wg.Wait()
}

func(pb *Publisher) publishEvent(event Event){
	if token := pb.Mqtt.Client.Subscribe(event.Topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())	}
	token := pb.Mqtt.Client.Publish(event.Topic, 0, false, event.Payload)
	token.Wait()
}



