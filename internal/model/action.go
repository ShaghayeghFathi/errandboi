package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Action struct {
	ID         primitive.ObjectID `bson:"_id"`
	Type       []string           `bson:"type"`
	EventCount int                `bson:"event_count"`
	Status     string             `bson:"status"`
}

type Event struct {
	ID          string      `bson:"_id"`
	ActionID    string      `bson:"action_id"`
	Description string      `bson:"description"`
	Delay       string      `bson:"delay"`
	Topic       string      `bson:"topic"`
	Payload     interface{} `json:"payload"`
	Status      string      `bson:"status"`
}
