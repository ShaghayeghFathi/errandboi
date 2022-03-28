package mongo

import (
	"context"
	"errandboi/internal/http/response"
	"errandboi/internal/model"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ActionCollection = "actions"
var EventsCollection = "events"

type MongoDB struct {
	DB *mongo.Database
}

func NewMongoDB(database *mongo.Database) *MongoDB {
	return &MongoDB{DB: database}
}

func (s *MongoDB) StoreEvent(ctx context.Context, event *model.Event) (string, error) {
	events := s.DB.Collection(EventsCollection)
	insertResult, err := events.InsertOne(ctx, event)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	return event.ID, nil
}

func (s *MongoDB) StoreAction(ctx context.Context, id primitive.ObjectID, t []string, eventCount int) (string, error) {
	actions := s.DB.Collection(ActionCollection)
	action := model.Action{
		ID:         id,
		Type:       t,
		EventCount: eventCount,
		Status:     "pending",
	}
	insertResult, err := actions.InsertOne(ctx, action)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document in actions: ", insertResult.InsertedID)

	return id.String(), nil
}

func (s *MongoDB) GetAction(ctx context.Context, id primitive.ObjectID) (model.Action, error) {
	res := s.DB.Collection(ActionCollection).FindOne(ctx, bson.M{
		"_id": id,
	})
	action := model.Action{}
	if err := res.Decode(&action); err != nil {
		fmt.Println(err)
	}
	return action, nil
}

func (s *MongoDB) GetEvents(ctx context.Context, actionId string) ([]response.EventResponse, error) {
	var events []response.EventResponse
	cursor, err := s.DB.Collection(EventsCollection).Find(ctx, bson.M{
		"action_id": actionId,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var res model.Event
		if err = cursor.Decode(&res); err != nil {
			log.Fatal(err)
		}
		result := response.EventResponse{Description: res.Description, Delay: res.Delay, Topic: res.Topic, Payload: res.Payload}
		events = append(events, result)
		fmt.Println(result)
	}
	return events, nil
}

func (s *MongoDB) GetEventStatus(ctx context.Context, actionId string) ([]response.EventStatusResponse, error) {
	var events []response.EventStatusResponse
	cursor, err := s.DB.Collection(EventsCollection).Find(ctx, bson.M{
		"action_id": actionId,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var res model.Event
		if err = cursor.Decode(&res); err != nil {
			log.Fatal(err)
		}
		result := response.EventStatusResponse{Description: res.Description, PublishDate: res.Delay, Status: res.Status}
		events = append(events, result)
	}
	return events, nil
}

func (s *MongoDB) GetEvent(ctx context.Context, eventId string, actionId string) (model.Event, error) {
	res := s.DB.Collection(EventsCollection).FindOne(ctx, bson.M{
		"_id":       eventId,
		"action_id": actionId,
	})
	event := model.Event{}
	if err := res.Decode(&event); err != nil {
		fmt.Println(err)
	}
	return event, nil
}
