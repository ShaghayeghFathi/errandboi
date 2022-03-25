package mongo

import (
	"context"
	"errandboi/internal/model"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ActionCollection = "actions"
var EventsCollection = "events"
type MongoDB struct {
	DB      *mongo.Database
}

func NewMongoDB(database *mongo.Database) *MongoDB{
	return &MongoDB{DB: database}
}

func (s *MongoDB) StoreEvent(ctx context.Context, id string, descp string, d string,topic string, payload interface{} ) (string, error) {
	events := s.DB.Collection(EventsCollection)
	event := model.Event{
		ID: id,
		Description: descp,
		Delay: d,
		Topic: topic,
		Payload: payload,
		Status: "pending",
	}
	insertResult, err := events.InsertOne(ctx, event)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	return id, nil
}

func (s *MongoDB) StoreAction(ctx context.Context, id primitive.ObjectID, t []string, eventCount int) (string, error) {
	actions := s.DB.Collection(ActionCollection)
	action := model.Action{
		ID: id,
		Type: t,
		EventCount: eventCount,
		Status: "pending",
	}
	insertResult, err := actions.InsertOne(ctx, action)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document in actions: ", insertResult.InsertedID)

	return id.String(), nil
}