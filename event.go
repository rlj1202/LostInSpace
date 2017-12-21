package lostinspace

type Event interface {
	Name() string
}

type EventListener interface {
	OnEvent(Event)
}

var eventlisteners []EventListener
var events []Event

func init() {
	eventlisteners = make([]EventListener, 0)
	events = make([]Event, 0, 10)
}

// Push event to queue.
func PushEvent(event Event) {
	events = append(events, event)
}

// Pop event from queue.
// Return nil if there is no more events remaining in queue.
func popEvent() Event {
	if len(events) == 0 {
		return nil
	}
	event := events[0]
	events = events[1:]

	return event
}

func RegisterEventListener(listener EventListener) {
	eventlisteners = append(eventlisteners, listener)
}

func PollEvents() {
	for event := popEvent(); event != nil; event = popEvent() {
		for _, listener := range eventlisteners {
			listener.OnEvent(event)
		}
	}
}
