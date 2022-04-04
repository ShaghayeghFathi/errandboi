package request

type Action struct {
	Type   []string `json:"type"`
	Events []Event  `json:"events"`
}

type Event struct {
	Description string      `json:"description"`
	Delay       string      `json:"delay"`
	Topic       string      `json:"topic"`
	Payload     interface{} `json:"payload"`
}
