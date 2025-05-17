package event

const (
	LinkVisited = "EVENT_LINK_VISITED"
)

type Event struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type EventBus struct {
	bus chan Event
}

func NewEventBus() *EventBus {
	return &EventBus{
		bus: make(chan Event),
	}
}

func (e *EventBus) Publish(event Event) {
	e.bus <- event
}

func (e *EventBus) Subscribe() <-chan Event {
	return e.bus
}
