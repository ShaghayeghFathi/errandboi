package response

import "time"

type GetEventsResponse struct {
	Type   []string        `json:"type"`
	Events []EventResponse `json:"events"`
}

type EventResponse struct {
	Description string      `json:"description"`
	Delay       string      `json:"delay"`
	Topic       string      `json:"topic"`
	Payload     interface{} `json:"payload"`
}

type GetEventsStatusResponse struct {
	Status string                `json:"status"`
	Events []EventStatusResponse `json:"events"`
}

type EventStatusResponse struct {
	Description string    `json:"description"`
	PublishDate time.Time `json:"publish_date"`
	Status      string    `json:"status"`
}
