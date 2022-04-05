package mongo

import (
	"context"
	"fmt"
	"log"

	"github.com/ShaghayeghFathi/errandboi/internal/http/response"
	"github.com/ShaghayeghFathi/errandboi/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ActionCollection = "actions"
	EventsCollection = "events"
)

type DB struct {
	DB *mongo.Database
}

func NewMongoDB(database *mongo.Database) *DB {
	return &DB{DB: database}
}

func (s *DB) StoreEvent(ctx context.Context, event *model.Event) (string, error) {
	events := s.DB.Collection(EventsCollection)

	_, err := events.InsertOne(ctx, event)
	if err != nil {
		log.Fatal(err)
	}

	return event.ID, nil
}

func (s *DB) StoreAction(ctx context.Context, id primitive.ObjectID, t []string, eventCount int) (string, error) {
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

func (s *DB) GetAction(ctx context.Context, id primitive.ObjectID) (model.Action, error) {
	res := s.DB.Collection(ActionCollection).FindOne(ctx, bson.M{
		"_id": id,
	})
	action := model.Action{}

	if err := res.Decode(&action); err != nil {
		fmt.Println(err)
	}

	return action, nil
}

func (s *DB) GetEvents(ctx context.Context, actionID string) ([]response.EventResponse, error) {
	var events []response.EventResponse

	cursor, err := s.DB.Collection(EventsCollection).Find(ctx, bson.M{
		"action_id": actionID,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var res model.Event
		if err = cursor.Decode(&res); err != nil {
			return nil, fmt.Errorf("cursor decode error %w", err)
		}

		result := response.EventResponse{
			Description: res.Description,
			Delay:       res.Delay,
			Topic:       res.Topic,
			Payload:     res.Payload,
		}

		events = append(events, result)
		fmt.Println(result)
	}

	return events, nil
}

func (s *DB) GetEventStatus(ctx context.Context, actionID string) ([]response.EventStatusResponse, error) {
	var events []response.EventStatusResponse

	cursor, err := s.DB.Collection(EventsCollection).Find(ctx, bson.M{
		"action_id": actionID,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var res model.Event
		if err = cursor.Decode(&res); err != nil {
			return nil, fmt.Errorf("cursor decode error %w", err)
		}

		result := response.EventStatusResponse{Description: res.Description, PublishDate: res.Delay, Status: res.Status}
		events = append(events, result)
	}

	return events, nil
}

func (s *DB) GetEvent(ctx context.Context, eventID string, actionID string) (model.Event, error) {
	res := s.DB.Collection(EventsCollection).FindOne(ctx, bson.M{
		"_id":       eventID,
		"action_id": actionID,
	})
	event := model.Event{}

	if err := res.Decode(&event); err != nil {
		fmt.Println(err)
	}

	return event, nil
}

func (s *DB) UpdateEventStatus(ctx context.Context, eventID string) *mongo.UpdateResult {
	res, err := s.DB.Collection(EventsCollection).UpdateOne(
		ctx,
		bson.M{"_id": eventID},
		bson.D{
			{Key: "$set", Value: bson.D{primitive.E{Key: "status", Value: "Done"}}},
		},
	)
	if err != nil {
		fmt.Println("update failed")
	}

	return res
}
