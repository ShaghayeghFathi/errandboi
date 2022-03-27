package publisher

import (
	"context"
	"errandboi/internal/model"
	"errandboi/internal/store/mongo"
	redisPK "errandboi/internal/store/redis"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gammazero/workerpool"
)

type Publisher struct{
	Redis *redisPK.RedisDB
	Mongo *mongo.MongoDB
	Events []model.Event
	wp         *workerpool.WorkerPool
	workerSize int
}

func NewPublisher(r *redisPK.RedisDB, m *mongo.MongoDB, size int) *Publisher{
	return &Publisher{Redis: r,Mongo: m, wp: workerpool.New(size), workerSize: size}
}

func(pb *Publisher) GetEvents(){
	var ctx = context.Background()
	start := float64(time.Now().Unix())
	events, err := pb.Redis.ZGetRange(ctx, "events", start, start+1)
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println("events from redis: ", events)
	for i := 0; i < len(events); i++ {
		eventId := events[i].Member.(string)
		fmt.Println("eventid is: ", eventId)
		s:= strings.Split(eventId, "_")
		event, err := pb.Mongo.GetEvent(context.Background(), eventId, s[0])
		if err!= nil{
			log.Fatal(err)
		}
		fmt.Println("event ", event)
		pb.Events = append(pb.Events, event)
	}
}

// func (pb *Publisher) Work() {
// 	var wg sync.WaitGroup

// 	for idx := range pb.Events {
// 		event := pb.Events[idx]
// 		wg.Add(1)
// 		pb.wp.Submit(func() {
// 			defer wg.Done()
// 			// publish event
// 		})
// 	}
// 	wg.Wait()
// }



