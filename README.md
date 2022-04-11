# Errandboi

Errandboi is a scheduler that publishes given events to Nats and EMQX. It recieves events via http request and caches them in Redis until their release time, then they get published to either EMQX or Nats or both. After events get published, their status changes to done.

## Instructions
### Lint and Run
* run docker containers, lint and run app: `make all`
### Build
* build app locally:  `make run build`
* build docker image locally: `make run docker-build`
### Run
* run services: `make up`
* run docker containers: `make run`

## Scheduling and publish structure

### scheduler
```go
func (sch *scheduler) WorkInIntervals(d time.Duration) {  
   ticker := time.NewTicker(d)  
  
   go func() {  
      for {  
         select {  
         case <-ticker.C:  
            sch.Publisher.GetEvents()  
            sch.Publisher.Work()  
         case <-sch.Stop:  
            sch.Publisher.Cancel()  
            sch.Publisher.Wp = workerpool.New(sch.Publisher.WorkerSize)  
  
            ticker.Stop()  
  
            return  
  }  
      }  
   }()  
}
```

### Nats
*creating stream

```go
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
```
* publish
```go
func (pb *Publisher) publishEventNats(event Event) {  
   t := natsp.ChannelName + "." + event.Topic  
   if _, err := pb.Nats.JSCtx.Publish(t, []byte(event.Payload)); err != nil {  
      pb.logger.Error("failed to publish event", zap.String("payload", event.Payload),  
  zap.String("topic", event.Topic), zap.Error(err))  
   }  
  
   pb.logger.Info("message published to nats", zap.String("payload", event.Payload))  
}
```
### EMQX
* publish
```go
func (pb *Publisher) publishEventEMQ(event Event) {  
   if token := pb.Mqtt.Client.Subscribe(event.Topic, 0, nil); token.Wait() && token.Error() != nil {  
      fmt.Println(token.Error())  
   }  
  
   token := pb.Mqtt.Client.Publish(event.Topic, 0, false, event.Payload)  
   pb.logger.Info("message published to emq: ", zap.String("payload", event.Payload))  
   token.Wait()  
}
```
## Example
#### Posting events
* request
```
curl --request POST '{{baseURL}}/events' \
--data-raw '{
	"type": [ "nats","emqx"],
	"events": [
		{
		"description": "This is an event description",
		"delay": "10s",
		"topic": "topic1",
		"payload": "hello"
		}
	]
}'
```
* response
```json
{
	"id": "62495362f8358c7b9d83ff4e"
}
```
#### Check event status
* request
```
curl --request GET '{{baseURL}}/events/{{eventId}}/status'
```
* response
```json
{
	"status": "done",
	"events": [
		{
		"description": "This is an event description",
		"publish_date":"2022-04-11T13:31:28+04:30",
		"status": "Done"
		}
	]
}
```