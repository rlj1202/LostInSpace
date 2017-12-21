package lostinspace

type KeyboardEvent struct {
	Key
	Scancode int
	Action
	Mods ModifierKey
}

type MouseEvent struct {
	XPos   float64
	YPos   float64
	Button MouseButton
	Action Action
	Mod    ModifierKey
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

func (event KeyboardEvent) Name() string {
	return "keyboardEvent"
}

func (event MouseEvent) Name() string {
	return "mouseEvent"
}

func (event ScrollEvent) Name() string {
	return "scrollEvent"
}

func (event CursorEnterEvent) Name() string {
	return "cursorEnterEvent"
}

func (event CursorPosEvent) Name() string {
	return "cursorPosEvent"
}
