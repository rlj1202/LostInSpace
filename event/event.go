package event

import "github.com/go-gl/glfw/v3.2/glfw"

type Event interface {
	GetEventName() string
}

type KeyboardEvent struct {
	Key      glfw.Key
	Scancode int
	Action   glfw.Action
	Mods     glfw.ModifierKey
}

type MouseEvent struct {
	XPos   float64
	YPos   float64
	Button glfw.MouseButton
	Action glfw.Action
	Mod    glfw.ModifierKey
}

type ScrollEvent struct {
	XOff float64
	YOff float64
}

type CursorEnterEvent struct {
	Entered bool
}

type CursorPosEvent struct {
	XPos float64
	YPos float64
}

var events []Event

func init() {
	events = make([]Event, 0, 10)
}

func (event KeyboardEvent) GetEventName() string {
	return "KeyboardEvent"
}

func (event MouseEvent) GetEventName() string {
	return "MouseEvent"
}

func (event ScrollEvent) GetEventName() string {
	return "ScrollEvent"
}

func (event CursorEnterEvent) GetEventName() string {
	return "CursorEnterEvent"
}

func (event CursorPosEvent) GetEventName() string {
	return "CursorPosEvent"
}

func Push(event Event) {
	events = append(events, event)
}

func Pop() Event {
	if len(events) == 0 {
		return nil
	}
	event := events[0]
	events = events[1:]

	return event
}
