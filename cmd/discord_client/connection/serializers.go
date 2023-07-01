package connection

type EventSerializer struct {
	Event interface{} `json:"event"`
	Type  string      `json:"type"`
}

func EventGetSerializerData(event interface{}, eventType string) EventSerializer {
	return EventSerializer{
		Event: event,
		Type:  eventType,
	}
}
